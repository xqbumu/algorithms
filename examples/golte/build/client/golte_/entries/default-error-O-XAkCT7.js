import { S as SvelteComponent, i as init, s as safe_not_equal, u as element, v as text, w as space, x as claim_element, y as children, z as claim_text, k as detach, A as claim_space, B as attr, d as insert_hydration, C as append_hydration, D as set_data, o as noop } from "../chunks/index-6zOx8Fy8.js";
function create_fragment(ctx) {
  let h1;
  let t0;
  let t1;
  let p;
  let t2;
  return {
    c() {
      h1 = element("h1");
      t0 = text(
        /*status*/
        ctx[0]
      );
      t1 = space();
      p = element("p");
      t2 = text(
        /*message*/
        ctx[1]
      );
      this.h();
    },
    l(nodes) {
      h1 = claim_element(nodes, "H1", {});
      var h1_nodes = children(h1);
      t0 = claim_text(
        h1_nodes,
        /*status*/
        ctx[0]
      );
      h1_nodes.forEach(detach);
      t1 = claim_space(nodes);
      p = claim_element(nodes, "P", { class: true });
      var p_nodes = children(p);
      t2 = claim_text(
        p_nodes,
        /*message*/
        ctx[1]
      );
      p_nodes.forEach(detach);
      this.h();
    },
    h() {
      attr(p, "class", "svelte-5zxcoy");
    },
    m(target, anchor) {
      insert_hydration(target, h1, anchor);
      append_hydration(h1, t0);
      insert_hydration(target, t1, anchor);
      insert_hydration(target, p, anchor);
      append_hydration(p, t2);
    },
    p(ctx2, [dirty]) {
      if (dirty & /*status*/
      1)
        set_data(
          t0,
          /*status*/
          ctx2[0]
        );
      if (dirty & /*message*/
      2)
        set_data(
          t2,
          /*message*/
          ctx2[1]
        );
    },
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
function instance($$self, $$props, $$invalidate) {
  let { status } = $$props;
  let { message } = $$props;
  $$self.$$set = ($$props2) => {
    if ("status" in $$props2)
      $$invalidate(0, status = $$props2.status);
    if ("message" in $$props2)
      $$invalidate(1, message = $$props2.message);
  };
  return [status, message];
}
class Default_error extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, instance, create_fragment, safe_not_equal, { status: 0, message: 1 });
  }
}
export {
  Default_error as default
};
//# sourceMappingURL=default-error-O-XAkCT7.js.map
