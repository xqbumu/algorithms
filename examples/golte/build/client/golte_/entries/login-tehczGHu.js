import { S as SvelteComponent, i as init, s as safe_not_equal, u as element, x as claim_element, G as get_svelte_dataset, B as attr, d as insert_hydration, o as noop, k as detach } from "../chunks/index-6zOx8Fy8.js";
function create_fragment(ctx) {
  let div;
  let textContent = `<input placeholder="Email Address" class="svelte-eww4au"/> <input placeholder="Password" type="password" class="svelte-eww4au"/> <button>Login</button>`;
  return {
    c() {
      div = element("div");
      div.innerHTML = textContent;
      this.h();
    },
    l(nodes) {
      div = claim_element(nodes, "DIV", { class: true, ["data-svelte-h"]: true });
      if (get_svelte_dataset(div) !== "svelte-6vslcr")
        div.innerHTML = textContent;
      this.h();
    },
    h() {
      attr(div, "class", "flex svelte-eww4au");
    },
    m(target, anchor) {
      insert_hydration(target, div, anchor);
    },
    p: noop,
    i: noop,
    o: noop,
    d(detaching) {
      if (detaching) {
        detach(div);
      }
    }
  };
}
class Login extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, null, create_fragment, safe_not_equal, {});
  }
}
export {
  Login as default
};
//# sourceMappingURL=login-tehczGHu.js.map
