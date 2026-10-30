// Run with: node --test openapi/rapidoc.test.mjs
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import test from "node:test";
import vm from "node:vm";

const source = readFileSync(new URL("./rapidoc.go", import.meta.url), "utf8");
const [pageScript, receiverScript] = [...source.matchAll(/<script>([\s\S]*?)<\/script>/g)].map((match) => match[1]);
const key = "gapi-rapidoc-oauth-request";
const callbackURL = "http://localhost:8080/oauth-receiver.html";

function storage(entries = []) {
  const values = new Map(entries);

  return {
    values,
    getItem(name) {
      return values.get(name) ?? null;
    },
    setItem(name, value) {
      values.set(name, value);
    },
    removeItem(name) {
      values.delete(name);
    },
  };
}

function loginWindow(state = "expected-state") {
  const parentStorage = storage();
  const microtasks = [];
  let beforeAuthorize;
  vm.runInNewContext(pageScript, {
    document: {
      querySelector() {
        return {
          addEventListener(name, listener) {
            assert.equal(name, "before-authorize");
            beforeAuthorize = listener;
          },
        };
      },
    },
    window: {
      location: new URL("http://localhost:8080/"),
      sessionStorage: parentStorage,
    },
    URL,
    queueMicrotask(task) {
      microtasks.push(task);
    },
  });
  beforeAuthorize({
    detail: {
      authParams: new URLSearchParams({ state, redirect_uri: callbackURL }),
    },
  });

  // A new window receives a copy; other windows opened after this task receive no pending login.
  const popupStorage = storage(parentStorage.values);

  for (const task of microtasks) {
    task();
  }

  assert.equal(parentStorage.getItem(key), null);

  return popupStorage;
}

function receive(suffix, popupStorage = storage(), supportsChannel = true) {
  const replies = [];
  const status = {
    textContent: "Completing authentication...",
  };
  const result = {
    replies,
    status,
    closed: false,
    channelClosed: false,
  };
  const window = {
    location: new URL(callbackURL + suffix),
    sessionStorage: popupStorage,
    history: {
      replaceState(_state, _title, path) {
        result.cleanedPath = path;
      },
    },
    close() {
      result.closed = true;
    },
  };
  Object.defineProperty(window, "opener", {
    get() {
      assert.equal(supportsChannel, false, "BroadcastChannel must avoid accessing the COOP-isolated opener");

      return {
        postMessage(reply, origin) {
          replies.push(JSON.parse(JSON.stringify(reply)));
          result.targetOrigin = origin;
        },
      };
    },
  });
  class Channel {
    constructor(name) {
      assert.equal(name, "rapidoc-oauth");
    }

    postMessage(reply) {
      replies.push(JSON.parse(JSON.stringify(reply)));
    }

    close() {
      result.channelClosed = true;
    }
  }
  vm.runInNewContext(receiverScript, {
    window,
    document: {
      getElementById() {
        return status;
      },
    },
    URLSearchParams,
    BroadcastChannel: supportsChannel ? Channel : undefined,
  });

  return result;
}

test("MWM Auth code without state completes from the popup's private login context", () => {
  const popupStorage = loginWindow();
  const result = receive("?code=test-code&expiry=1800000000&expire_in=3600", popupStorage);
  assert.equal(result.replies.length, 1);
  assert.equal(result.replies[0].state, "expected-state");
  assert.equal(result.replies[0].code, "test-code");
  assert.equal(result.cleanedPath, "/oauth-receiver.html");
  assert.equal(popupStorage.getItem(key), null);
  assert.equal(result.channelClosed, true);
  assert.equal(result.closed, true);
});

test("standard OAuth preserves a supplied state and implicit token replies", () => {
  const code = receive("?code=test-code&state=expected-state", loginWindow());
  assert.equal(code.replies[0].state, "expected-state");
  assert.equal(code.closed, true);
  const token = receive("#access_token=test-token&token_type=Bearer&state=expected-state", loginWindow());
  assert.equal(token.replies[0].responseType, "token");
  assert.equal(token.replies[0].accessToken, "test-token");
  assert.equal(token.closed, true);
});

test("a missing state is never recovered from an unrelated login window", () => {
  const first = loginWindow("first-state");
  const second = loginWindow("second-state");
  assert.equal(receive("?code=first-code", first).replies[0].state, "first-state");
  assert.equal(receive("?code=second-code", second).replies[0].state, "second-state");
  const unsolicited = receive("?code=untrusted-code");
  assert.equal(unsolicited.replies.length, 0);
  assert.equal(unsolicited.closed, false);
  assert.match(unsolicited.status.textContent, /missing or expired/);
});

test("a conflicting state is rejected without replacing it", () => {
  const result = receive("?code=test-code&state=other-state", loginWindow());
  assert.equal(result.replies.length, 0);
  assert.equal(result.closed, false);
  assert.match(result.status.textContent, /does not match/);
});

test("expired or different callback contexts cannot restore a missing state", () => {
  for (const change of [
    (request) => {
      request.startedAt -= 16 * 60 * 1000;
    },
    (request) => {
      request.startedAt += 60 * 1000;
    },
    (request) => {
      request.callback = "https://other.example/oauth-receiver.html";
    },
  ]) {
    const popupStorage = loginWindow();
    const request = JSON.parse(popupStorage.getItem(key));
    change(request);
    popupStorage.setItem(key, JSON.stringify(request));
    const result = receive("?code=test-code", popupStorage);
    assert.equal(result.replies.length, 0);
    assert.equal(result.closed, false);
  }
});

test("provider errors return to the UI, and empty callbacks show a useful message", () => {
  const error = receive("?error=access_denied&error_description=Login+cancelled", loginWindow());
  assert.equal(error.replies[0].error, "access_denied");
  assert.equal(error.replies[0].errorDescription, "Login cancelled");
  assert.equal(error.closed, true);
  const empty = receive("");
  assert.equal(empty.replies.length, 0);
  assert.equal(empty.closed, false);
  assert.match(empty.status.textContent, /No OAuth response/);
});

test("legacy browsers use postMessage with the exact callback origin", () => {
  const result = receive("?code=test-code", loginWindow(), false);
  assert.equal(result.replies[0].state, "expected-state");
  assert.equal(result.targetOrigin, "http://localhost:8080");
  assert.equal(result.closed, true);
});
