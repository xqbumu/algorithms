var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => {
  __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
  return value;
};
import { o as noop, s as safe_not_equal, r as get_store_value } from "./index-6zOx8Fy8.js";
const scriptRel = "modulepreload";
const assetsURL = function(dep) {
  return "/" + dep;
};
const seen = {};
const __vitePreload = function preload(baseModule, deps, importerUrl) {
  let promise = Promise.resolve();
  if (deps && deps.length > 0) {
    const links = document.getElementsByTagName("link");
    promise = Promise.all(deps.map((dep) => {
      dep = assetsURL(dep);
      if (dep in seen)
        return;
      seen[dep] = true;
      const isCss = dep.endsWith(".css");
      const cssSelector = isCss ? '[rel="stylesheet"]' : "";
      const isBaseRelative = !!importerUrl;
      if (isBaseRelative) {
        for (let i = links.length - 1; i >= 0; i--) {
          const link2 = links[i];
          if (link2.href === dep && (!isCss || link2.rel === "stylesheet")) {
            return;
          }
        }
      } else if (document.querySelector(`link[href="${dep}"]${cssSelector}`)) {
        return;
      }
      const link = document.createElement("link");
      link.rel = isCss ? "stylesheet" : scriptRel;
      if (!isCss) {
        link.as = "script";
        link.crossOrigin = "";
      }
      link.href = dep;
      document.head.appendChild(link);
      if (isCss) {
        return new Promise((res, rej) => {
          link.addEventListener("load", res);
          link.addEventListener("error", () => rej(new Error(`Unable to preload CSS for ${dep}`)));
        });
      }
    }));
  }
  return promise.then(() => baseModule()).catch((err) => {
    const e = new Event("vite:preloadError", { cancelable: true });
    e.payload = err;
    window.dispatchEvent(e);
    if (!e.defaultPrevented) {
      throw err;
    }
  });
};
const golteContext = Symbol();
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
  function subscribe(run, invalidate = noop) {
    const subscriber = [run, invalidate];
    subscribers.add(subscriber);
    if (subscribers.size === 1) {
      stop = start(set, update) || noop;
    }
    run(value);
    return () => {
      subscribers.delete(subscriber);
      if (subscribers.size === 0 && stop) {
        stop();
        stop = null;
      }
    };
  }
  return { set, update, subscribe };
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
  constructor(url, nodes) {
    __publicField(this, "url");
    __publicField(this, "node");
    this.url = writable(new URL(url));
    this.node = fromArray(nodes);
  }
}
class ClientAppState extends ServerAppState {
  constructor(url, nodes) {
    super(url, nodes);
    __publicField(this, "hrefMap", {});
    this.hrefMap[get_store_value(this.url).href] = new Promise((r) => r(nodes));
  }
  async update(href) {
    this.url.set(new URL(href));
    const array = await (this.hrefMap[href] ?? load(href));
    let before = this.node;
    let after = fromArray(array);
    while (true) {
      const bval = get_store_value(before);
      const aval = get_store_value(after);
      if (!bval && !aval)
        break;
      const bcomp = bval == null ? void 0 : bval.content.comp;
      const acomp = aval == null ? void 0 : aval.content.comp;
      if (bcomp === acomp) {
        before = bval.next;
        after = aval.next;
      } else {
        before.set(aval);
        break;
      }
    }
  }
}
const AppState = ClientAppState;
async function load(href) {
  const headers = { "Golte": "true" };
  const resp = await fetch(href, { headers });
  const json = await resp.json();
  for (const entry of [...json.Entries, json.ErrPage]) {
    for (const css of entry.CSS) {
      if (document.querySelector(`link[href="${css}"][rel="stylesheet"]`))
        continue;
      const link = document.createElement("link");
      link.href = css;
      link.rel = "stylesheet";
      document.head.appendChild(link);
    }
  }
  const promises = json.Entries.map(async (entry) => ({
    comp: (await __vitePreload(() => import(entry.File), true ? __vite__mapDeps([]) : void 0)).default,
    props: entry.Props,
    errPage: (await __vitePreload(() => import(json.ErrPage.File), true ? __vite__mapDeps([]) : void 0)).default
  }));
  return await Promise.all(promises);
}
export {
  AppState as A,
  __vitePreload as _,
  golteContext as g,
  load as l
};
function __vite__mapDeps(indexes) {
  if (!__vite__mapDeps.viteFileDeps) {
    __vite__mapDeps.viteFileDeps = []
  }
  return indexes.map((i) => __vite__mapDeps.viteFileDeps[i])
}
//# sourceMappingURL=appstate-tLgYl6N1.js.map
