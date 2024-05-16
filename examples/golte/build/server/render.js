"use strict";
var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => {
  __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
  return value;
};
Object.defineProperty(exports, Symbol.toStringTag, { value: "Module" });
function noop() {
}
function run(fn) {
  return fn();
}
function blank_object() {
  return /* @__PURE__ */ Object.create(null);
}
function run_all(fns) {
  fns.forEach(run);
}
function safe_not_equal(a, b) {
  return a != a ? b == b : a !== b || a && typeof a === "object" || typeof a === "function";
}
function subscribe(store, ...callbacks) {
  if (store == null) {
    for (const callback of callbacks) {
      callback(void 0);
    }
    return noop;
  }
  const unsub = store.subscribe(...callbacks);
  return unsub.unsubscribe ? () => unsub.unsubscribe() : unsub;
}
let current_component;
function set_current_component(component) {
  current_component = component;
}
function get_current_component() {
  if (!current_component)
    throw new Error("Function called outside component initialization");
  return current_component;
}
function setContext(key, context) {
  get_current_component().$$.context.set(key, context);
  return context;
}
function getContext(key) {
  return get_current_component().$$.context.get(key);
}
const ATTR_REGEX = /[&"]/g;
const CONTENT_REGEX = /[&<]/g;
function escape(value, is_attr = false) {
  const str = String(value);
  const pattern = is_attr ? ATTR_REGEX : CONTENT_REGEX;
  pattern.lastIndex = 0;
  let escaped = "";
  let last = 0;
  while (pattern.test(str)) {
    const i = pattern.lastIndex - 1;
    const ch = str[i];
    escaped += str.substring(last, i) + (ch === "&" ? "&amp;" : ch === '"' ? "&quot;" : "&lt;");
    last = i + 1;
  }
  return escaped + str.substring(last);
}
const missing_component = {
  $$render: () => ""
};
function validate_component(component, name) {
  if (!component || !component.$$render) {
    if (name === "svelte:component")
      name += " this={...}";
    throw new Error(
      `<${name}> is not a valid SSR component. You may need to review your build config to ensure that dependencies are compiled, rather than imported as pre-compiled modules. Otherwise you may need to fix a <${name}>.`
    );
  }
  return component;
}
let on_destroy;
function create_ssr_component(fn) {
  function $$render(result, props, bindings, slots, context) {
    const parent_component = current_component;
    const $$ = {
      on_destroy,
      context: new Map(context || (parent_component ? parent_component.$$.context : [])),
      // these will be immediately discarded
      on_mount: [],
      before_update: [],
      after_update: [],
      callbacks: blank_object()
    };
    set_current_component({ $$ });
    const html = fn(result, props, bindings, slots);
    set_current_component(parent_component);
    return html;
  }
  return {
    render: (props = {}, { $$slots = {}, context = /* @__PURE__ */ new Map() } = {}) => {
      on_destroy = [];
      const result = { title: "", head: "", css: /* @__PURE__ */ new Set() };
      const html = $$render(result, props, {}, $$slots, context);
      run_all(on_destroy);
      return {
        html,
        css: {
          code: Array.from(result.css).map((css2) => css2.code).join("\n"),
          map: null
          // TODO
        },
        head: result.title + result.head
      };
    },
    $$render
  };
}
function add_attribute(name, value, boolean) {
  if (value == null || boolean && !value)
    return "";
  const assignment = boolean && value === true ? "" : `="${escape(value, true)}"`;
  return ` ${name}${assignment}`;
}
const Node_1 = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $next, $$unsubscribe_next;
  let { node } = $$props;
  let { index } = $$props;
  const { next, content } = node;
  $$unsubscribe_next = subscribe(next, (value) => $next = value);
  if ($$props.node === void 0 && $$bindings.node && node !== void 0)
    $$bindings.node(node);
  if ($$props.index === void 0 && $$bindings.index && index !== void 0)
    $$bindings.index(index);
  $$unsubscribe_next();
  return `  ${validate_component(content.comp || missing_component, "svelte:component").$$render($$result, Object.assign({}, content.props), {}, {
    default: () => {
      return ` ${$next ? ` ${validate_component(Node, "Node").$$render($$result, { node: $next, index: index + 1 }, {}, {})}` : ``}`;
    }
  })}`;
});
const golteContext = Symbol();
const handleError = Symbol();
const ServerNode = Node_1;
const ssrWrapper = {
  ...ServerNode,
  $$render: (result, props, bindings, slots, context) => {
    try {
      return ServerNode.$$render(result, props, bindings, slots, context);
    } catch (err) {
      let message = "Internal Error";
      {
        message = err instanceof Error && err.stack ? err.stack : String(err);
      }
      const errProps = {
        status: 500,
        message
      };
      getContext(handleError)({ index: props.index, props: errProps });
      return props.node.content.errPage.$$render(result, errProps, bindings, slots, context);
    }
  }
};
const Node = ssrWrapper;
const subscriber_queue = [];
function writable(value, start = noop) {
  let stop;
  const subscribers = /* @__PURE__ */ new Set();
  function set(new_value) {
    if (safe_not_equal(value, new_value)) {
      value = new_value;
      if (stop) {
        const run_queue = !subscriber_queue.length;
        for (const subscriber of subscribers) {
          subscriber[1]();
          subscriber_queue.push(subscriber, value);
        }
        if (run_queue) {
          for (let i = 0; i < subscriber_queue.length; i += 2) {
            subscriber_queue[i][0](subscriber_queue[i + 1]);
          }
          subscriber_queue.length = 0;
        }
      }
    }
  }
  function update(fn) {
    set(fn(value));
  }
  function subscribe2(run2, invalidate = noop) {
    const subscriber = [run2, invalidate];
    subscribers.add(subscriber);
    if (subscribers.size === 1) {
      stop = start(set, update) || noop;
    }
    run2(value);
    return () => {
      subscribers.delete(subscriber);
      if (subscribers.size === 0 && stop) {
        stop();
        stop = null;
      }
    };
  }
  return { set, update, subscribe: subscribe2 };
}
function fromArray(array) {
  let current = writable(null);
  for (let i = array.length - 1; i >= 0; i--) {
    current = writable({
      content: array[i],
      next: current
    });
  }
  return current;
}
class ServerAppState {
  constructor(url2, nodes) {
    __publicField(this, "url");
    __publicField(this, "node");
    this.url = writable(new URL(url2));
    this.node = fromArray(nodes);
  }
}
const AppState = ServerAppState;
const Root$1 = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $node, $$unsubscribe_node;
  let { nodes } = $$props;
  let { contextData } = $$props;
  const state = new AppState(contextData.URL, nodes);
  const { node } = state;
  $$unsubscribe_node = subscribe(node, (value) => $node = value);
  setContext(golteContext, state);
  if ($$props.nodes === void 0 && $$bindings.nodes && nodes !== void 0)
    $$bindings.nodes(nodes);
  if ($$props.contextData === void 0 && $$bindings.contextData && contextData !== void 0)
    $$bindings.contextData(contextData);
  $$unsubscribe_node();
  return ` ${$node ? `${validate_component(Node, "Node").$$render($$result, { node: $node, index: 0 }, {}, {})}` : ``}`;
});
function getGolteContext() {
  return getContext(golteContext);
}
const url = {
  subscribe(fn) {
    return getGolteContext().url.subscribe(fn);
  }
};
const css$4 = {
  code: 'main.svelte-auera{padding:20px}ul.svelte-auera{margin:0;height:60px;list-style:none;background:#ccddf8;;;display:flex;gap:10px;justify-content:center;align-items:center}li.svelte-auera{height:40px;padding:0 10px}li[aria-current="page"].svelte-auera{border-bottom:solid 2px orangered}a.svelte-auera{height:100%;display:flex;align-items:center;color:inherit;font-family:"Gill Sans", sans-serif;;;font-weight:700;font-size:0.9rem;text-decoration:none}a.svelte-auera:hover{color:orangered}',
  map: null
};
const Main = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $url, $$unsubscribe_url;
  $$unsubscribe_url = subscribe(url, (value) => $url = value);
  $$result.css.add(css$4);
  $$unsubscribe_url();
  return `<header><nav><ul class="svelte-auera"><li${add_attribute("aria-current", $url.pathname === "/" ? "page" : void 0, 0)} class="svelte-auera"><a href="/" class="svelte-auera" data-svelte-h="svelte-1s2d22s">Home</a></li> <li${add_attribute("aria-current", $url.pathname === "/about" ? "page" : void 0, 0)} class="svelte-auera"><a href="/about" class="svelte-auera" data-svelte-h="svelte-fjk9up">About</a></li> <li${add_attribute("aria-current", $url.pathname === "/contact" ? "page" : void 0, 0)} class="svelte-auera"><a href="/contact" class="svelte-auera" data-svelte-h="svelte-9ha5c7">Contact</a></li> <li${add_attribute("aria-current", $url.pathname === "/user/profile" ? "page" : void 0, 0)} class="svelte-auera"><a href="/user/profile" class="svelte-auera" data-svelte-h="svelte-11kmhqb">Profile</a></li> <li${add_attribute("aria-current", $url.pathname === "/user/login" ? "page" : void 0, 0)} class="svelte-auera"><a href="/user/login" class="svelte-auera" data-svelte-h="svelte-j77ggn">Login</a></li></ul></nav></header> <main class="svelte-auera">${slots.default ? slots.default({}) : ``} </main>`;
});
const css$3 = {
  code: "div.svelte-cyctq1{margin:0 auto;width:60%;padding:20px;border:1px solid #333333;border-radius:12px}",
  map: null
};
const Secondary = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  $$result.css.add(css$3);
  return `<div class="svelte-cyctq1">${slots.default ? slots.default({}) : ``} </div>`;
});
const About = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<h1 data-svelte-h="svelte-1ihvjf0">Sample About Page</h1> <p data-svelte-h="svelte-rzm3sy">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>`;
});
const Contact = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<h1 data-svelte-h="svelte-hq2lh">Sample Contact Page</h1> <p data-svelte-h="svelte-rzm3sy">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>`;
});
const Home = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<h1 data-svelte-h="svelte-zr09k0">Sample Home Page</h1> <p data-svelte-h="svelte-rzm3sy">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>`;
});
const css$2 = {
  code: ".flex.svelte-eww4au{display:flex;flex-direction:column;gap:20px}input.svelte-eww4au{height:50px}",
  map: null
};
const Login = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  $$result.css.add(css$2);
  return `<div class="flex svelte-eww4au" data-svelte-h="svelte-6vslcr"><input placeholder="Email Address" class="svelte-eww4au"> <input placeholder="Password" type="password" class="svelte-eww4au"> <button>Login</button> </div>`;
});
const css$1 = {
  code: "h2.svelte-19p5cuz.svelte-19p5cuz{margin-bottom:0}.name.svelte-19p5cuz.svelte-19p5cuz{padding-bottom:20px}tr.svelte-19p5cuz td.svelte-19p5cuz:first-of-type{font-weight:700;padding-right:20px}",
  map: null
};
const Profile = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { username } = $$props;
  let { realname } = $$props;
  let { occupation } = $$props;
  let { age } = $$props;
  let { email } = $$props;
  let { site } = $$props;
  let { searching } = $$props;
  if ($$props.username === void 0 && $$bindings.username && username !== void 0)
    $$bindings.username(username);
  if ($$props.realname === void 0 && $$bindings.realname && realname !== void 0)
    $$bindings.realname(realname);
  if ($$props.occupation === void 0 && $$bindings.occupation && occupation !== void 0)
    $$bindings.occupation(occupation);
  if ($$props.age === void 0 && $$bindings.age && age !== void 0)
    $$bindings.age(age);
  if ($$props.email === void 0 && $$bindings.email && email !== void 0)
    $$bindings.email(email);
  if ($$props.site === void 0 && $$bindings.site && site !== void 0)
    $$bindings.site(site);
  if ($$props.searching === void 0 && $$bindings.searching && searching !== void 0)
    $$bindings.searching(searching);
  $$result.css.add(css$1);
  return `<div class="name svelte-19p5cuz"><h2 class="svelte-19p5cuz">${escape(realname)}</h2> <small>${escape(occupation)}</small></div> <table><tbody><tr class="svelte-19p5cuz"><td class="svelte-19p5cuz" data-svelte-h="svelte-1ccubqy">Username</td> <td class="svelte-19p5cuz">${escape(username)}</td></tr> <tr class="svelte-19p5cuz"><td class="svelte-19p5cuz" data-svelte-h="svelte-1t5ooc7">Age</td> <td class="svelte-19p5cuz">${escape(age)}</td></tr> <tr class="svelte-19p5cuz"><td class="svelte-19p5cuz" data-svelte-h="svelte-s939w6">Email</td> <td class="svelte-19p5cuz">${escape(email)}</td></tr> <tr class="svelte-19p5cuz"><td class="svelte-19p5cuz" data-svelte-h="svelte-2cfzwp">Site</td> <td class="svelte-19p5cuz"><a${add_attribute("href", site, 0)}>${escape(site)}</a></td></tr> <tr class="svelte-19p5cuz"><td class="svelte-19p5cuz" data-svelte-h="svelte-8hox4k">Looking for job?</td> <td class="svelte-19p5cuz">${escape(searching ? "Yes" : "No")}</td></tr></tbody> </table>`;
});
const css = {
  code: "p.svelte-5zxcoy{white-space:pre-wrap}",
  map: null
};
const Default_error = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { status } = $$props;
  let { message } = $$props;
  if ($$props.status === void 0 && $$bindings.status && status !== void 0)
    $$bindings.status(status);
  if ($$props.message === void 0 && $$bindings.message && message !== void 0)
    $$bindings.message(message);
  $$result.css.add(css);
  return `<h1>${escape(status)}</h1> <p class="svelte-5zxcoy">${escape(message)}</p>`;
});
const Root = Root$1;
const hydrate = "/golte_/entries/hydrate--REBL6VH.js";
const Manifest = {
  "layout/main": {
    server: Main,
    Client: "/golte_/entries/main-v9Tf3z2U.js",
    CSS: [
      "/golte_/assets/main-A4ha7JCj.css"
    ]
  },
  "layout/secondary": {
    server: Secondary,
    Client: "/golte_/entries/secondary-lVDqtvFS.js",
    CSS: [
      "/golte_/assets/secondary-d_XPN1dB.css"
    ]
  },
  "page/about": {
    server: About,
    Client: "/golte_/entries/about-wEaqGju3.js",
    CSS: []
  },
  "page/contact": {
    server: Contact,
    Client: "/golte_/entries/contact-4qyY_YF-.js",
    CSS: []
  },
  "page/home": {
    server: Home,
    Client: "/golte_/entries/home-TFhQMUgX.js",
    CSS: []
  },
  "page/login": {
    server: Login,
    Client: "/golte_/entries/login-tehczGHu.js",
    CSS: [
      "/golte_/assets/login-dXZNShM1.css"
    ]
  },
  "page/profile": {
    server: Profile,
    Client: "/golte_/entries/profile-KyvgK15P.js",
    CSS: [
      "/golte_/assets/profile-jyOHKTMr.css"
    ]
  },
  "$$$GOLTE_DEFAULT_ERROR$$$": {
    server: Default_error,
    Client: "/golte_/entries/default-error-O-XAkCT7.js",
    CSS: [
      "/golte_/assets/default-error-B8dlbLlS.css"
    ]
  }
};
function Render(entries, contextData, errPage) {
  const serverNodes = [];
  const clientNodes = [];
  const stylesheets = /* @__PURE__ */ new Set();
  const err = Manifest[errPage];
  if (!err)
    throw new Error(`"${errPage}" is not a component`);
  for (const e of entries) {
    const c = Manifest[e.Comp];
    if (!c)
      throw new Error(`"${e.Comp}" is not a component`);
    serverNodes.push({ comp: c.server, props: e.Props, errPage: err.server });
    clientNodes.push({ comp: `${c.Client}`, props: e.Props, errPage: `${err.Client}` });
    for (const path of c.CSS) {
      stylesheets.add(path);
    }
  }
  for (const path of err.CSS) {
    stylesheets.add(path);
  }
  let error;
  const context = /* @__PURE__ */ new Map();
  context.set(handleError, (e) => error = e);
  let { html, head } = Root.render({ nodes: serverNodes, contextData }, { context });
  for (const path of stylesheets) {
    head += `
<link href="${path}" rel="stylesheet">`;
  }
  if (error) {
    clientNodes[error.index].ssrError = error.props;
  }
  html += `
        <script>
            (async function () {
                const target = document.currentScript.parentElement;
                const { hydrate } = await import("${hydrate}");
                await hydrate(target, ${stringify(clientNodes)}, ${stringify(contextData)});
            })();
        <\/script>
    `;
  return {
    Head: head,
    Body: html,
    HasError: !!error
  };
}
function stringify(object) {
  return JSON.stringify(object).replace("<\/script>", "<\\/script>");
}
exports.Manifest = Manifest;
exports.Render = Render;
//# sourceMappingURL=render.js.map
