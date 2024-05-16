import { S as SvelteComponent, i as init, s as safe_not_equal, u as element, w as space, x as claim_element, G as get_svelte_dataset, A as claim_space, d as insert_hydration, o as noop, k as detach } from "../chunks/index-6zOx8Fy8.js";
function create_fragment(ctx) {
  let h1;
  let textContent = "Sample About Page";
  let t1;
  let p;
  let textContent_1 = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.";
  return {
    c() {
      h1 = element("h1");
      h1.textContent = textContent;
      t1 = space();
      p = element("p");
      p.textContent = textContent_1;
    },
    l(nodes) {
      h1 = claim_element(nodes, "H1", { ["data-svelte-h"]: true });
      if (get_svelte_dataset(h1) !== "svelte-1ihvjf0")
        h1.textContent = textContent;
      t1 = claim_space(nodes);
      p = claim_element(nodes, "P", { ["data-svelte-h"]: true });
      if (get_svelte_dataset(p) !== "svelte-rzm3sy")
        p.textContent = textContent_1;
    },
    m(target, anchor) {
      insert_hydration(target, h1, anchor);
      insert_hydration(target, t1, anchor);
      insert_hydration(target, p, anchor);
    },
    p: noop,
    i: noop,
    o: noop,
    d(detaching) {
      if (detaching) {
        detach(h1);
        detach(t1);
        detach(p);
      }
    }
  };
}
class About extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, null, create_fragment, safe_not_equal, {});
  }
}
export {
  About as default
};
//# sourceMappingURL=about-wEaqGju3.js.map
