package mcp

import (
	"net/http"
	"strings"

	"github.com/mwm-io/gapi/handler"
)

// InspectorHandler serves the MCP Inspector HTML page.
type InspectorHandler struct {
	handler.WithMiddlewares
	mcpEndpoint string
	serverName  string
}

// NewInspectorHandler builds a new InspectorHandler.
func NewInspectorHandler(config Config, middlewares ...handler.Middleware) handler.Handler {
	return &InspectorHandler{
		mcpEndpoint: config.GetPath(),
		serverName:  config.GetServerName("MCP Server"),
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: middlewares,
		},
	}
}

// Serve implements the handler.Handler interface.
func (h InspectorHandler) Serve(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return strings.NewReader(inspectorHTML(h.mcpEndpoint, h.serverName)), nil
}

func inspectorHTML(mcpEndpoint, serverName string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>MCP Inspector — ` + serverName + `</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --font-sans:"Inter",system-ui,-apple-system,sans-serif;
  --font-mono:"JetBrains Mono","SF Mono",Monaco,Consolas,monospace;
  --bg-root:#0a0f1a;
  --bg-surface:#111827;
  --bg-elevated:#1a2332;
  --bg-input:#0d1320;
  --bg-hover:rgba(99,91,255,.06);
  --bg-active:rgba(99,91,255,.1);
  --text-primary:#f1f5f9;
  --text-secondary:#8b95a8;
  --text-tertiary:#5a6578;
  --text-accent:#a78bfa;
  --border:#1e293b;
  --border-focus:#635bff;
  --primary:#635bff;
  --primary-hover:#7c75ff;
  --primary-glow:rgba(99,91,255,.15);
  --success:#10b981;
  --success-bg:rgba(16,185,129,.08);
  --success-border:rgba(16,185,129,.2);
  --error:#ef4444;
  --error-bg:rgba(239,68,68,.08);
  --error-border:rgba(239,68,68,.2);
  --warning:#f59e0b;
  --radius-sm:6px;
  --radius-md:8px;
  --radius-lg:12px;
  --shadow-sm:0 1px 2px rgba(0,0,0,.3);
  --shadow-md:0 4px 12px rgba(0,0,0,.25);
  --shadow-lg:0 8px 24px rgba(0,0,0,.3);
  --transition:150ms cubic-bezier(.4,0,.2,1);
}
html,body{height:100%;font-family:var(--font-sans);background:var(--bg-root);color:var(--text-primary);font-size:14px;line-height:1.5;-webkit-font-smoothing:antialiased}

/* ── Layout ── */
#app{display:flex;flex-direction:column;height:100vh}

/* ── Header ── */
header{
  background:var(--bg-surface);border-bottom:1px solid var(--border);
  padding:0 24px;height:56px;display:flex;align-items:center;gap:16px;flex-shrink:0;
}
.logo{display:flex;align-items:center;gap:10px}
.logo svg{width:20px;height:20px;color:var(--primary)}
.logo h1{font-size:14px;font-weight:600;color:var(--text-primary);letter-spacing:-.01em}
.sep{width:1px;height:20px;background:var(--border)}
.server-badge{
  font-size:12px;font-weight:500;color:var(--text-secondary);
  background:var(--bg-elevated);border:1px solid var(--border);
  padding:3px 10px;border-radius:20px;
}
.status-pill{
  margin-left:auto;display:flex;align-items:center;gap:6px;
  font-size:12px;font-weight:500;color:var(--text-tertiary);
}
.status-dot{width:7px;height:7px;border-radius:50%;background:var(--text-tertiary);flex-shrink:0}
.status-pill.connected .status-dot{background:var(--success);box-shadow:0 0 6px var(--success)}
.status-pill.connected{color:var(--text-secondary)}
.status-pill.error .status-dot{background:var(--error);box-shadow:0 0 6px var(--error)}
.status-pill.error{color:var(--error)}

/* ── Main layout ── */
.main{display:flex;flex:1;overflow:hidden}

/* ── Sidebar ── */
.sidebar{
  width:280px;min-width:240px;background:var(--bg-surface);
  border-right:1px solid var(--border);display:flex;flex-direction:column;flex-shrink:0;
}
.sidebar-top{padding:16px 16px 12px;border-bottom:1px solid var(--border)}
.sidebar-title{font-size:11px;font-weight:600;color:var(--text-tertiary);text-transform:uppercase;letter-spacing:.08em;margin-bottom:10px}
.filter-wrap{position:relative}
.filter-wrap svg{position:absolute;left:10px;top:50%;transform:translateY(-50%);width:14px;height:14px;color:var(--text-tertiary);pointer-events:none}
.filter-input{
  width:100%;background:var(--bg-input);border:1px solid var(--border);
  border-radius:var(--radius-sm);padding:7px 10px 7px 32px;
  color:var(--text-primary);font-size:13px;font-family:var(--font-sans);outline:none;
  transition:border-color var(--transition),box-shadow var(--transition);
}
.filter-input:focus{border-color:var(--border-focus);box-shadow:0 0 0 3px var(--primary-glow)}
.filter-input::placeholder{color:var(--text-tertiary)}
.tool-list{overflow-y:auto;flex:1;padding:6px 8px}
.tool-list::-webkit-scrollbar{width:4px}
.tool-list::-webkit-scrollbar-track{background:transparent}
.tool-list::-webkit-scrollbar-thumb{background:var(--border);border-radius:4px}
.tool-item{
  padding:10px 12px;cursor:pointer;border-radius:var(--radius-sm);
  transition:background var(--transition);margin-bottom:2px;
}
.tool-item:hover{background:var(--bg-hover)}
.tool-item.active{background:var(--bg-active)}
.tool-item .t-name{
  font-size:13px;font-weight:500;color:var(--text-primary);
  font-family:var(--font-mono);letter-spacing:-.01em;
}
.tool-item.active .t-name{color:var(--primary)}
.tool-item .t-desc{
  font-size:12px;color:var(--text-tertiary);margin-top:2px;
  white-space:nowrap;overflow:hidden;text-overflow:ellipsis;line-height:1.4;
}
.tool-count{padding:12px 16px;border-top:1px solid var(--border);font-size:11px;color:var(--text-tertiary);text-align:center}

/* ── Content ── */
.content{flex:1;overflow-y:auto;background:var(--bg-root)}
.content::-webkit-scrollbar{width:6px}
.content::-webkit-scrollbar-track{background:transparent}
.content::-webkit-scrollbar-thumb{background:var(--border);border-radius:4px}
.content-inner{max-width:720px;margin:0 auto;padding:40px 48px}

/* ── Empty state ── */
.empty-state{
  display:flex;flex-direction:column;align-items:center;justify-content:center;
  height:100%;gap:16px;color:var(--text-tertiary);
}
.empty-state svg{width:48px;height:48px;opacity:.3}
.empty-state p{font-size:14px}
.empty-state .hint{font-size:12px;color:var(--text-tertiary);opacity:.6}

/* ── Tool detail ── */
.tool-header{margin-bottom:32px}
.tool-name{
  font-size:24px;font-weight:600;color:var(--text-primary);
  font-family:var(--font-mono);letter-spacing:-.02em;
}
.tool-desc{color:var(--text-secondary);margin-top:8px;line-height:1.6;font-size:14px}

/* ── Section ── */
.section-label{
  font-size:11px;font-weight:600;color:var(--text-tertiary);
  text-transform:uppercase;letter-spacing:.08em;margin-bottom:16px;
  padding-bottom:8px;border-bottom:1px solid var(--border);
}

/* ── Param rows ── */
.params-list{display:flex;flex-direction:column;gap:1px;background:var(--border);border-radius:var(--radius-md);overflow:hidden;box-shadow:var(--shadow-sm)}
.param-row{
  display:flex;align-items:flex-start;gap:16px;
  padding:14px 16px;background:var(--bg-surface);
  transition:background var(--transition);
}
.param-row:hover{background:var(--bg-elevated)}
.param-info{flex:0 0 200px}
.param-name{font-size:13px;font-weight:500;color:var(--text-primary);font-family:var(--font-mono)}
.param-badges{display:flex;gap:6px;margin-top:4px;flex-wrap:wrap}
.badge{
  display:inline-block;font-size:10px;font-weight:600;
  padding:1px 6px;border-radius:4px;letter-spacing:.02em;
}
.badge-type{color:var(--text-accent);background:rgba(167,139,250,.1)}
.badge-required{color:var(--error);background:var(--error-bg)}
.badge-optional{color:var(--text-tertiary);background:rgba(90,101,120,.15)}
.param-desc{font-size:12px;color:var(--text-tertiary);margin-top:4px;line-height:1.4}
.param-input{flex:1;min-width:0}
.param-input input,.param-input textarea,.param-input select{
  width:100%;background:var(--bg-input);border:1px solid var(--border);
  border-radius:var(--radius-sm);padding:8px 12px;
  color:var(--text-primary);font-size:13px;font-family:var(--font-sans);outline:none;
  transition:border-color var(--transition),box-shadow var(--transition);
}
.param-input textarea{
  min-height:72px;resize:vertical;font-family:var(--font-mono);font-size:12px;line-height:1.5;
}
.param-input input:focus,.param-input textarea:focus,.param-input select:focus{
  border-color:var(--border-focus);box-shadow:0 0 0 3px var(--primary-glow);
}
.param-input input::placeholder,.param-input textarea::placeholder{color:var(--text-tertiary)}
.param-input select{
  cursor:pointer;-webkit-appearance:none;appearance:none;
  background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='none' viewBox='0 0 12 12'%3E%3Cpath stroke='%235a6578' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round' d='M3 4.5 6 7.5 9 4.5'/%3E%3C/svg%3E");
  background-repeat:no-repeat;background-position:right 12px center;padding-right:32px;
}
.param-input select option{background:var(--bg-surface);color:var(--text-primary)}
.no-params{color:var(--text-tertiary);font-size:13px;padding:16px;background:var(--bg-surface);border-radius:var(--radius-md);border:1px solid var(--border)}

/* ── Actions ── */
.actions{margin-top:24px;display:flex;gap:10px;align-items:center}
.btn{
  display:inline-flex;align-items:center;gap:6px;
  padding:9px 20px;border:none;border-radius:var(--radius-sm);
  font-size:13px;font-weight:500;font-family:var(--font-sans);
  cursor:pointer;transition:all var(--transition);outline:none;
}
.btn-primary{
  background:var(--primary);color:#fff;box-shadow:0 1px 2px rgba(99,91,255,.3);
}
.btn-primary:hover{background:var(--primary-hover);box-shadow:0 2px 8px rgba(99,91,255,.35)}
.btn-primary:active{transform:scale(.98)}
.btn-primary:disabled{opacity:.5;cursor:not-allowed;transform:none}
.btn-ghost{background:transparent;color:var(--text-secondary);padding:9px 14px}
.btn-ghost:hover{color:var(--text-primary);background:var(--bg-hover)}
.btn svg{width:14px;height:14px}
.timer{font-size:12px;color:var(--text-tertiary);font-family:var(--font-mono)}

/* ── Response ── */
.response{margin-top:32px;animation:fadeIn .2s ease}
@keyframes fadeIn{from{opacity:0;transform:translateY(4px)}to{opacity:1;transform:translateY(0)}}
.response-header{display:flex;align-items:center;gap:10px;margin-bottom:10px}
.response-label{font-size:11px;font-weight:600;color:var(--text-tertiary);text-transform:uppercase;letter-spacing:.08em}
.response-status{font-size:11px;font-weight:600;padding:2px 8px;border-radius:4px}
.response-status.ok{color:var(--success);background:var(--success-bg)}
.response-status.err{color:var(--error);background:var(--error-bg)}
.response-box{
  background:var(--bg-surface);border:1px solid var(--border);
  border-radius:var(--radius-md);overflow:hidden;box-shadow:var(--shadow-sm);
}
.response-box.is-error{border-color:var(--error-border)}
.response-box.is-success{border-color:var(--success-border)}
.response-code{
  padding:16px 20px;font-family:var(--font-mono);font-size:12px;
  line-height:1.7;white-space:pre-wrap;word-break:break-word;
  max-height:420px;overflow-y:auto;color:var(--text-secondary);
}
.response-code::-webkit-scrollbar{width:4px}
.response-code::-webkit-scrollbar-track{background:transparent}
.response-code::-webkit-scrollbar-thumb{background:var(--border);border-radius:4px}

/* ── JSON syntax ── */
.json-key{color:#a78bfa}
.json-str{color:#34d399}
.json-num{color:#fbbf24}
.json-bool{color:#f472b6}
.json-null{color:#64748b}

/* ── Spinner ── */
.spinner{
  width:14px;height:14px;border:2px solid rgba(255,255,255,.2);
  border-top-color:#fff;border-radius:50%;animation:spin .5s linear infinite;
}
@keyframes spin{to{transform:rotate(360deg)}}

/* ── Loading skeleton ── */
.skeleton-list{padding:6px 8px}
.skeleton-item{padding:10px 12px;margin-bottom:2px}
.skeleton-line{height:12px;background:var(--bg-elevated);border-radius:4px;animation:pulse 1.5s ease-in-out infinite}
.skeleton-line.w60{width:60%}
.skeleton-line.w80{width:80%}
.skeleton-line.w40{width:40%;margin-top:6px}
@keyframes pulse{0%,100%{opacity:.4}50%{opacity:.8}}
</style>
</head>
<body>
<div id="app">
  <header>
    <div class="logo">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/></svg>
      <h1>MCP Inspector</h1>
    </div>
    <div class="sep"></div>
    <span class="server-badge" id="server-name">` + serverName + `</span>
    <div class="status-pill" id="status"><span class="status-dot"></span><span id="status-text">Connecting</span></div>
  </header>
  <div class="main">
    <div class="sidebar">
      <div class="sidebar-top">
        <div class="sidebar-title">Tools</div>
        <div class="filter-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
          <input type="text" class="filter-input" id="filter" placeholder="Search tools...">
        </div>
      </div>
      <div class="tool-list" id="tool-list">
        <div class="skeleton-list">
          <div class="skeleton-item"><div class="skeleton-line w80"></div><div class="skeleton-line w40"></div></div>
          <div class="skeleton-item"><div class="skeleton-line w60"></div><div class="skeleton-line w40"></div></div>
          <div class="skeleton-item"><div class="skeleton-line w80"></div><div class="skeleton-line w40"></div></div>
        </div>
      </div>
      <div class="tool-count" id="tool-count"></div>
    </div>
    <div class="content" id="content">
      <div class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
        <p>Select a tool to get started</p>
        <span class="hint">Choose from the sidebar or use the search to filter</span>
      </div>
    </div>
  </div>
</div>
<script>
(function(){
  const MCP_ENDPOINT = "` + mcpEndpoint + `";
  let sessionId = null;
  let reqId = 0;
  let tools = [];
  let selectedTool = null;

  const $ = s => document.querySelector(s);
  const statusEl = $("#status");
  const statusTextEl = $("#status-text");
  const toolListEl = $("#tool-list");
  const toolCountEl = $("#tool-count");
  const contentEl = $("#content");
  const filterEl = $("#filter");

  function setStatus(text, cls) {
    statusTextEl.textContent = text;
    statusEl.className = "status-pill" + (cls ? " " + cls : "");
  }

  async function rpc(method, params) {
    const id = ++reqId;
    const body = { jsonrpc: "2.0", id, method };
    if (params) body.params = params;
    const headers = { "Content-Type": "application/json", "Accept": "application/json, text/event-stream" };
    if (sessionId) headers["Mcp-Session-Id"] = sessionId;
    const resp = await fetch(MCP_ENDPOINT, { method: "POST", headers, body: JSON.stringify(body) });
    const sid = resp.headers.get("Mcp-Session-Id");
    if (sid) sessionId = sid;
    const ct = resp.headers.get("Content-Type") || "";
    if (ct.includes("text/event-stream")) return parseSSE(await resp.text(), id);
    return resp.json();
  }

  async function notify(method, params) {
    const body = { jsonrpc: "2.0", method };
    if (params) body.params = params;
    const headers = { "Content-Type": "application/json" };
    if (sessionId) headers["Mcp-Session-Id"] = sessionId;
    await fetch(MCP_ENDPOINT, { method: "POST", headers, body: JSON.stringify(body) });
  }

  function parseSSE(text, expectedId) {
    for (const line of text.split("\n")) {
      if (!line.startsWith("data: ")) continue;
      try {
        const msg = JSON.parse(line.slice(6));
        if (msg.id === expectedId || msg.result || msg.error) return msg;
      } catch(e) {}
    }
    return null;
  }

  async function initialize() {
    try {
      setStatus("Initializing...", "");
      const initResult = await rpc("initialize", {
        protocolVersion: "2025-03-26",
        capabilities: {},
        clientInfo: { name: "GAPI MCP Inspector", version: "1.0.0" }
      });
      if (initResult?.result?.serverInfo) {
        const info = initResult.result.serverInfo;
        $("#server-name").textContent = info.name + " v" + info.version;
      }
      await notify("notifications/initialized");
      const toolsResult = await rpc("tools/list");
      if (toolsResult?.result?.tools) {
        tools = toolsResult.result.tools;
        tools.sort((a, b) => a.name.localeCompare(b.name));
      }
      setStatus("Connected", "connected");
      toolCountEl.textContent = tools.length + " tool" + (tools.length !== 1 ? "s" : "") + " available";
      renderToolList();
    } catch(e) {
      setStatus("Failed", "error");
      toolListEl.innerHTML = "";
      toolCountEl.textContent = "Connection error";
    }
  }

  function renderToolList(filter) {
    const f = (filter || "").toLowerCase();
    const filtered = f ? tools.filter(t => t.name.toLowerCase().includes(f) || (t.description || "").toLowerCase().includes(f)) : tools;
    toolListEl.innerHTML = "";
    if (filtered.length === 0 && tools.length > 0) {
      toolListEl.innerHTML = '<div style="padding:24px 16px;text-align:center;color:var(--text-tertiary);font-size:13px">No tools match your search</div>';
      return;
    }
    filtered.forEach(tool => {
      const div = document.createElement("div");
      div.className = "tool-item" + (selectedTool?.name === tool.name ? " active" : "");
      div.innerHTML = '<div class="t-name">' + esc(tool.name) + '</div>' +
        (tool.description ? '<div class="t-desc">' + esc(tool.description) + '</div>' : '');
      div.onclick = () => selectTool(tool);
      toolListEl.appendChild(div);
    });
  }

  function selectTool(tool) {
    selectedTool = tool;
    renderToolList(filterEl.value);
    renderToolDetail(tool);
  }

  function renderToolDetail(tool) {
    const schema = tool.inputSchema || {};
    const props = schema.properties || {};
    const required = schema.required || [];
    const propNames = Object.keys(props);

    let html = '<div class="content-inner"><div class="tool-header">';
    html += '<div class="tool-name">' + esc(tool.name) + '</div>';
    if (tool.description) html += '<div class="tool-desc">' + esc(tool.description) + '</div>';
    html += '</div>';

    html += '<div class="section-label">Parameters</div>';
    if (propNames.length === 0) {
      html += '<div class="no-params">This tool takes no parameters</div>';
    } else {
      html += '<div class="params-list">';
      propNames.forEach(name => {
        const prop = props[name];
        const isReq = required.includes(name);
        const type = prop.type || "string";
        html += '<div class="param-row"><div class="param-info">';
        html += '<div class="param-name">' + esc(name) + '</div>';
        html += '<div class="param-badges">';
        html += '<span class="badge badge-type">' + esc(type) + '</span>';
        html += isReq ? '<span class="badge badge-required">required</span>' : '<span class="badge badge-optional">optional</span>';
        html += '</div>';
        if (prop.description) html += '<div class="param-desc">' + esc(prop.description) + '</div>';
        html += '</div><div class="param-input">' + renderInput(name, prop, type) + '</div></div>';
      });
      html += '</div>';
    }

    html += '<div class="actions">';
    html += '<button class="btn btn-primary" id="call-btn"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>Execute</button>';
    html += '<button class="btn btn-ghost" id="clear-btn">Clear</button>';
    html += '<span class="timer" id="timer"></span>';
    html += '</div>';
    html += '<div id="response-area"></div>';
    html += '</div>';

    contentEl.innerHTML = html;
    $("#call-btn").onclick = () => callTool(tool, propNames, props);
    $("#clear-btn").onclick = () => clearForm(propNames);
  }

  function renderInput(name, prop, type) {
    const attr = 'data-param="' + esc(name) + '"';
    if (prop.enum?.length > 0) {
      let s = '<select ' + attr + '><option value="">Select a value</option>';
      prop.enum.forEach(v => { s += '<option value="' + esc(String(v)) + '">' + esc(String(v)) + '</option>'; });
      return s + '</select>';
    }
    if (type === "boolean") return '<select ' + attr + '><option value="">Select</option><option value="true">true</option><option value="false">false</option></select>';
    if (type === "object" || type === "array") return '<textarea ' + attr + ' placeholder=\'{"key": "value"}\'></textarea>';
    const inputType = (type === "integer" || type === "number") ? "number" : "text";
    const step = type === "number" ? ' step="any"' : '';
    const ph = type === "integer" ? "0" : type === "number" ? "0.0" : "Enter value";
    return '<input type="' + inputType + '"' + step + ' ' + attr + ' placeholder="' + ph + '">';
  }

  function clearForm(propNames) {
    propNames.forEach(name => {
      const el = document.querySelector('[data-param="' + name + '"]');
      if (el) el.value = "";
    });
    const ra = $("#response-area");
    if (ra) ra.innerHTML = "";
  }

  async function callTool(tool, propNames, props) {
    const args = {};
    for (const name of propNames) {
      const el = document.querySelector('[data-param="' + name + '"]');
      if (!el || el.value === "") continue;
      const type = (props[name].type || "string");
      if (type === "integer") args[name] = parseInt(el.value, 10);
      else if (type === "number") args[name] = parseFloat(el.value);
      else if (type === "boolean") args[name] = el.value === "true";
      else if (type === "object" || type === "array") {
        try { args[name] = JSON.parse(el.value); } catch(e) { showResponse("Invalid JSON for " + name + ": " + e.message, true, 0); return; }
      } else args[name] = el.value;
    }

    const btn = $("#call-btn");
    const timerEl = $("#timer");
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span>Executing...';
    const t0 = performance.now();

    try {
      const result = await rpc("tools/call", { name: tool.name, arguments: args });
      const ms = Math.round(performance.now() - t0);
      if (result?.error) {
        showResponse(JSON.stringify(result.error, null, 2), true, ms);
      } else if (result?.result) {
        const isErr = result.result.isError === true;
        const content = result.result.content || [];
        const text = content.map(c => c.type === "text" ? c.text : JSON.stringify(c)).join("\n");
        showResponse(tryFormatJSON(text), isErr, ms);
      } else {
        showResponse("No response received", true, Math.round(performance.now() - t0));
      }
    } catch(e) {
      showResponse("Request failed: " + e.message, true, Math.round(performance.now() - t0));
    } finally {
      btn.disabled = false;
      btn.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>Execute';
    }
  }

  function showResponse(text, isError, ms) {
    const area = $("#response-area");
    const statusLabel = isError ? '<span class="response-status err">Error</span>' : '<span class="response-status ok">Success</span>';
    const timeLabel = ms > 0 ? '<span class="timer">' + ms + 'ms</span>' : '';
    area.innerHTML = '<div class="response">' +
      '<div class="response-header"><span class="response-label">Response</span>' + statusLabel + timeLabel + '</div>' +
      '<div class="response-box ' + (isError ? 'is-error' : 'is-success') + '">' +
      '<div class="response-code">' + syntaxHighlight(text) + '</div></div></div>';
  }

  function tryFormatJSON(text) {
    try { return JSON.stringify(JSON.parse(text), null, 2); } catch(e) { return text; }
  }

  function syntaxHighlight(json) {
    return esc(json).replace(
      /("(?:\\.|[^"\\])*")\s*:/g,
      '<span class="json-key">$1</span>:'
    ).replace(
      /:\s*("(?:\\.|[^"\\])*")/g,
      ': <span class="json-str">$1</span>'
    ).replace(
      /:\s*(-?\d+\.?\d*(?:[eE][+-]?\d+)?)/g,
      ': <span class="json-num">$1</span>'
    ).replace(
      /:\s*(true|false)/g,
      ': <span class="json-bool">$1</span>'
    ).replace(
      /:\s*(null)/g,
      ': <span class="json-null">$1</span>'
    );
  }

  function esc(s) {
    const d = document.createElement("div");
    d.textContent = s;
    return d.innerHTML;
  }

  filterEl.addEventListener("input", () => renderToolList(filterEl.value));
  initialize();
})();
</script>
</body>
</html>`
}
