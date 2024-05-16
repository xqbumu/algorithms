import { E as getContext, S as SvelteComponent, i as init, s as safe_not_equal, F as create_slot, u as element, w as space, x as claim_element, y as children, G as get_svelte_dataset, k as detach, A as claim_space, B as attr, d as insert_hydration, C as append_hydration, H as action_destroyer, I as update_slot_base, J as get_all_dirty_from_scope, K as get_slot_changes, j as transition_in, t as transition_out, L as run_all, l as component_subscribe } from "../chunks/index-6zOx8Fy8.js";
import { g as golteContext, l as load } from "../chunks/appstate-tLgYl6N1.js";
const preload = (a, preload2 = "hover") => {
  const state = getContext(golteContext);
  async function loadAnchor() {
    if (a.origin !== location.origin)
      return;
    if (a.href in state.hrefMap)
      return;
    state.hrefMap[a.href] = load(a.href);
  }
  if (preload2 === "mount")
    loadAnchor();
  if (preload2 === "hover")
    a.addEventListener("mouseover", loadAnchor);
  if (preload2 === "tap") {
    a.addEventListener("mousedown", loadAnchor);
    a.addEventListener("touchstart", loadAnchor);
  }
  a.addEventListener("click", async (e) => {
    if (a.origin !== location.origin)
      return;
    e.preventDefault();
    await state.update(a.href);
    history.pushState(a.href, "", a.href);
  });
};
function getGolteContext() {
  return getContext(golteContext);
}
const url = {
  subscribe(fn) {
    return getGolteContext().url.subscribe(fn);
  }
};
function create_fragment(ctx) {
  let header;
  let nav;
  let ul;
  let li0;
  let a0;
  let textContent = "Home";
  let li0_aria_current_value;
  let t1;
  let li1;
  let a1;
  let textContent_1 = "About";
  let li1_aria_current_value;
  let t3;
  let li2;
  let a2;
  let textContent_2 = "Contact";
  let li2_aria_current_value;
  let t5;
  let li3;
  let a3;
  let textContent_3 = "Profile";
  let li3_aria_current_value;
  let t7;
  let li4;
  let a4;
  let textContent_4 = "Login";
  let li4_aria_current_value;
  let t9;
  let main;
  let current;
  let mounted;
  let dispose;
  const default_slot_template = (
    /*#slots*/
    ctx[2].default
  );
  const default_slot = create_slot(
    default_slot_template,
    ctx,
    /*$$scope*/
    ctx[1],
    null
  );
  return {
    c() {
      header = element("header");
      nav = element("nav");
      ul = element("ul");
      li0 = element("li");
      a0 = element("a");
      a0.textContent = textContent;
      t1 = space();
      li1 = element("li");
      a1 = element("a");
      a1.textContent = textContent_1;
      t3 = space();
      li2 = element("li");
      a2 = element("a");
      a2.textContent = textContent_2;
      t5 = space();
      li3 = element("li");
      a3 = element("a");
      a3.textContent = textContent_3;
      t7 = space();
      li4 = element("li");
      a4 = element("a");
      a4.textContent = textContent_4;
      t9 = space();
      main = element("main");
      if (default_slot)
        default_slot.c();
      this.h();
    },
    l(nodes) {
      header = claim_element(nodes, "HEADER", {});
      var header_nodes = children(header);
      nav = claim_element(header_nodes, "NAV", {});
      var nav_nodes = children(nav);
      ul = claim_element(nav_nodes, "UL", { class: true });
      var ul_nodes = children(ul);
      li0 = claim_element(ul_nodes, "LI", { "aria-current": true, class: true });
      var li0_nodes = children(li0);
      a0 = claim_element(li0_nodes, "A", {
        href: true,
        class: true,
        ["data-svelte-h"]: true
      });
      if (get_svelte_dataset(a0) !== "svelte-1s2d22s")
        a0.textContent = textContent;
      li0_nodes.forEach(detach);
      t1 = claim_space(ul_nodes);
      li1 = claim_element(ul_nodes, "LI", { "aria-current": true, class: true });
      var li1_nodes = children(li1);
      a1 = claim_element(li1_nodes, "A", {
        href: true,
        class: true,
        ["data-svelte-h"]: true
      });
      if (get_svelte_dataset(a1) !== "svelte-fjk9up")
        a1.textContent = textContent_1;
      li1_nodes.forEach(detach);
      t3 = claim_space(ul_nodes);
      li2 = claim_element(ul_nodes, "LI", { "aria-current": true, class: true });
      var li2_nodes = children(li2);
      a2 = claim_element(li2_nodes, "A", {
        href: true,
        class: true,
        ["data-svelte-h"]: true
      });
      if (get_svelte_dataset(a2) !== "svelte-9ha5c7")
        a2.textContent = textContent_2;
      li2_nodes.forEach(detach);
      t5 = claim_space(ul_nodes);
      li3 = claim_element(ul_nodes, "LI", { "aria-current": true, class: true });
      var li3_nodes = children(li3);
      a3 = claim_element(li3_nodes, "A", {
        href: true,
        class: true,
        ["data-svelte-h"]: true
      });
      if (get_svelte_dataset(a3) !== "svelte-11kmhqb")
        a3.textContent = textContent_3;
      li3_nodes.forEach(detach);
      t7 = claim_space(ul_nodes);
      li4 = claim_element(ul_nodes, "LI", { "aria-current": true, class: true });
      var li4_nodes = children(li4);
      a4 = claim_element(li4_nodes, "A", {
        href: true,
        class: true,
        ["data-svelte-h"]: true
      });
      if (get_svelte_dataset(a4) !== "svelte-j77ggn")
        a4.textContent = textContent_4;
      li4_nodes.forEach(detach);
      ul_nodes.forEach(detach);
      nav_nodes.forEach(detach);
      header_nodes.forEach(detach);
      t9 = claim_space(nodes);
      main = claim_element(nodes, "MAIN", { class: true });
      var main_nodes = children(main);
      if (default_slot)
        default_slot.l(main_nodes);
      main_nodes.forEach(detach);
      this.h();
    },
    h() {
      attr(a0, "href", "/");
      attr(a0, "class", "svelte-auera");
      attr(li0, "aria-current", li0_aria_current_value = /*$url*/
      ctx[0].pathname === "/" ? "page" : void 0);
      attr(li0, "class", "svelte-auera");
      attr(a1, "href", "/about");
      attr(a1, "class", "svelte-auera");
      attr(li1, "aria-current", li1_aria_current_value = /*$url*/
      ctx[0].pathname === "/about" ? "page" : void 0);
      attr(li1, "class", "svelte-auera");
      attr(a2, "href", "/contact");
      attr(a2, "class", "svelte-auera");
      attr(li2, "aria-current", li2_aria_current_value = /*$url*/
      ctx[0].pathname === "/contact" ? "page" : void 0);
      attr(li2, "class", "svelte-auera");
      attr(a3, "href", "/user/profile");
      attr(a3, "class", "svelte-auera");
      attr(li3, "aria-current", li3_aria_current_value = /*$url*/
      ctx[0].pathname === "/user/profile" ? "page" : void 0);
      attr(li3, "class", "svelte-auera");
      attr(a4, "href", "/user/login");
      attr(a4, "class", "svelte-auera");
      attr(li4, "aria-current", li4_aria_current_value = /*$url*/
      ctx[0].pathname === "/user/login" ? "page" : void 0);
      attr(li4, "class", "svelte-auera");
      attr(ul, "class", "svelte-auera");
      attr(main, "class", "svelte-auera");
    },
    m(target, anchor) {
      insert_hydration(target, header, anchor);
      append_hydration(header, nav);
      append_hydration(nav, ul);
      append_hydration(ul, li0);
      append_hydration(li0, a0);
      append_hydration(ul, t1);
      append_hydration(ul, li1);
      append_hydration(li1, a1);
      append_hydration(ul, t3);
      append_hydration(ul, li2);
      append_hydration(li2, a2);
      append_hydration(ul, t5);
      append_hydration(ul, li3);
      append_hydration(li3, a3);
      append_hydration(ul, t7);
      append_hydration(ul, li4);
      append_hydration(li4, a4);
      insert_hydration(target, t9, anchor);
      insert_hydration(target, main, anchor);
      if (default_slot) {
        default_slot.m(main, null);
      }
      current = true;
      if (!mounted) {
        dispose = [
          action_destroyer(preload.call(null, a0)),
          action_destroyer(preload.call(null, a1)),
          action_destroyer(preload.call(null, a2)),
          action_destroyer(preload.call(null, a3)),
          action_destroyer(preload.call(null, a4))
        ];
        mounted = true;
      }
    },
    p(ctx2, [dirty]) {
      if (!current || dirty & /*$url*/
      1 && li0_aria_current_value !== (li0_aria_current_value = /*$url*/
      ctx2[0].pathname === "/" ? "page" : void 0)) {
        attr(li0, "aria-current", li0_aria_current_value);
      }
      if (!current || dirty & /*$url*/
      1 && li1_aria_current_value !== (li1_aria_current_value = /*$url*/
      ctx2[0].pathname === "/about" ? "page" : void 0)) {
        attr(li1, "aria-current", li1_aria_current_value);
      }
      if (!current || dirty & /*$url*/
      1 && li2_aria_current_value !== (li2_aria_current_value = /*$url*/
      ctx2[0].pathname === "/contact" ? "page" : void 0)) {
        attr(li2, "aria-current", li2_aria_current_value);
      }
      if (!current || dirty & /*$url*/
      1 && li3_aria_current_value !== (li3_aria_current_value = /*$url*/
      ctx2[0].pathname === "/user/profile" ? "page" : void 0)) {
        attr(li3, "aria-current", li3_aria_current_value);
      }
      if (!current || dirty & /*$url*/
      1 && li4_aria_current_value !== (li4_aria_current_value = /*$url*/
      ctx2[0].pathname === "/user/login" ? "page" : void 0)) {
        attr(li4, "aria-current", li4_aria_current_value);
      }
      if (default_slot) {
        if (default_slot.p && (!current || dirty & /*$$scope*/
        2)) {
          update_slot_base(
            default_slot,
            default_slot_template,
            ctx2,
            /*$$scope*/
            ctx2[1],
            !current ? get_all_dirty_from_scope(
              /*$$scope*/
              ctx2[1]
            ) : get_slot_changes(
              default_slot_template,
              /*$$scope*/
              ctx2[1],
              dirty,
              null
            ),
            null
          );
        }
      }
    },
    i(local) {
      if (current)
        return;
      transition_in(default_slot, local);
      current = true;
    },
    o(local) {
      transition_out(default_slot, local);
      current = false;
    },
    d(detaching) {
      if (detaching) {
        detach(header);
        detach(t9);
        detach(main);
      }
      if (default_slot)
        default_slot.d(detaching);
      mounted = false;
      run_all(dispose);
    }
  };
}
function instance($$self, $$props, $$invalidate) {
  let $url;
  component_subscribe($$self, url, ($$value) => $$invalidate(0, $url = $$value));
  let { $$slots: slots = {}, $$scope } = $$props;
  $$self.$$set = ($$props2) => {
    if ("$$scope" in $$props2)
      $$invalidate(1, $$scope = $$props2.$$scope);
  };
  return [$url, $$scope, slots];
}
class Main extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, instance, create_fragment, safe_not_equal, {});
  }
}
export {
  Main as default
};
//# sourceMappingURL=main-v9Tf3z2U.js.map
