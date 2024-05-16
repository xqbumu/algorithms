import { S as SvelteComponent, i as init, s as safe_not_equal, u as element, v as text, w as space, x as claim_element, y as children, z as claim_text, k as detach, A as claim_space, G as get_svelte_dataset, B as attr, d as insert_hydration, C as append_hydration, D as set_data, o as noop } from "../chunks/index-6zOx8Fy8.js";
function create_fragment(ctx) {
  let div;
  let h2;
  let t0;
  let t1;
  let small;
  let t2;
  let t3;
  let table;
  let tbody;
  let tr0;
  let td0;
  let textContent = "Username";
  let t5;
  let td1;
  let t6;
  let t7;
  let tr1;
  let td2;
  let textContent_1 = "Age";
  let t9;
  let td3;
  let t10;
  let t11;
  let tr2;
  let td4;
  let textContent_2 = "Email";
  let t13;
  let td5;
  let t14;
  let t15;
  let tr3;
  let td6;
  let textContent_3 = "Site";
  let t17;
  let td7;
  let a;
  let t18;
  let t19;
  let tr4;
  let td8;
  let textContent_4 = "Looking for job?";
  let t21;
  let td9;
  let t22_value = (
    /*searching*/
    ctx[6] ? "Yes" : "No"
  );
  let t22;
  return {
    c() {
      div = element("div");
      h2 = element("h2");
      t0 = text(
        /*realname*/
        ctx[1]
      );
      t1 = space();
      small = element("small");
      t2 = text(
        /*occupation*/
        ctx[2]
      );
      t3 = space();
      table = element("table");
      tbody = element("tbody");
      tr0 = element("tr");
      td0 = element("td");
      td0.textContent = textContent;
      t5 = space();
      td1 = element("td");
      t6 = text(
        /*username*/
        ctx[0]
      );
      t7 = space();
      tr1 = element("tr");
      td2 = element("td");
      td2.textContent = textContent_1;
      t9 = space();
      td3 = element("td");
      t10 = text(
        /*age*/
        ctx[3]
      );
      t11 = space();
      tr2 = element("tr");
      td4 = element("td");
      td4.textContent = textContent_2;
      t13 = space();
      td5 = element("td");
      t14 = text(
        /*email*/
        ctx[4]
      );
      t15 = space();
      tr3 = element("tr");
      td6 = element("td");
      td6.textContent = textContent_3;
      t17 = space();
      td7 = element("td");
      a = element("a");
      t18 = text(
        /*site*/
        ctx[5]
      );
      t19 = space();
      tr4 = element("tr");
      td8 = element("td");
      td8.textContent = textContent_4;
      t21 = space();
      td9 = element("td");
      t22 = text(t22_value);
      this.h();
    },
    l(nodes) {
      div = claim_element(nodes, "DIV", { class: true });
      var div_nodes = children(div);
      h2 = claim_element(div_nodes, "H2", { class: true });
      var h2_nodes = children(h2);
      t0 = claim_text(
        h2_nodes,
        /*realname*/
        ctx[1]
      );
      h2_nodes.forEach(detach);
      t1 = claim_space(div_nodes);
      small = claim_element(div_nodes, "SMALL", {});
      var small_nodes = children(small);
      t2 = claim_text(
        small_nodes,
        /*occupation*/
        ctx[2]
      );
      small_nodes.forEach(detach);
      div_nodes.forEach(detach);
      t3 = claim_space(nodes);
      table = claim_element(nodes, "TABLE", {});
      var table_nodes = children(table);
      tbody = claim_element(table_nodes, "TBODY", {});
      var tbody_nodes = children(tbody);
      tr0 = claim_element(tbody_nodes, "TR", { class: true });
      var tr0_nodes = children(tr0);
      td0 = claim_element(tr0_nodes, "TD", { class: true, ["data-svelte-h"]: true });
      if (get_svelte_dataset(td0) !== "svelte-1ccubqy")
        td0.textContent = textContent;
      t5 = claim_space(tr0_nodes);
      td1 = claim_element(tr0_nodes, "TD", { class: true });
      var td1_nodes = children(td1);
      t6 = claim_text(
        td1_nodes,
        /*username*/
        ctx[0]
      );
      td1_nodes.forEach(detach);
      tr0_nodes.forEach(detach);
      t7 = claim_space(tbody_nodes);
      tr1 = claim_element(tbody_nodes, "TR", { class: true });
      var tr1_nodes = children(tr1);
      td2 = claim_element(tr1_nodes, "TD", { class: true, ["data-svelte-h"]: true });
      if (get_svelte_dataset(td2) !== "svelte-1t5ooc7")
        td2.textContent = textContent_1;
      t9 = claim_space(tr1_nodes);
      td3 = claim_element(tr1_nodes, "TD", { class: true });
      var td3_nodes = children(td3);
      t10 = claim_text(
        td3_nodes,
        /*age*/
        ctx[3]
      );
      td3_nodes.forEach(detach);
      tr1_nodes.forEach(detach);
      t11 = claim_space(tbody_nodes);
      tr2 = claim_element(tbody_nodes, "TR", { class: true });
      var tr2_nodes = children(tr2);
      td4 = claim_element(tr2_nodes, "TD", { class: true, ["data-svelte-h"]: true });
      if (get_svelte_dataset(td4) !== "svelte-s939w6")
        td4.textContent = textContent_2;
      t13 = claim_space(tr2_nodes);
      td5 = claim_element(tr2_nodes, "TD", { class: true });
      var td5_nodes = children(td5);
      t14 = claim_text(
        td5_nodes,
        /*email*/
        ctx[4]
      );
      td5_nodes.forEach(detach);
      tr2_nodes.forEach(detach);
      t15 = claim_space(tbody_nodes);
      tr3 = claim_element(tbody_nodes, "TR", { class: true });
      var tr3_nodes = children(tr3);
      td6 = claim_element(tr3_nodes, "TD", { class: true, ["data-svelte-h"]: true });
      if (get_svelte_dataset(td6) !== "svelte-2cfzwp")
        td6.textContent = textContent_3;
      t17 = claim_space(tr3_nodes);
      td7 = claim_element(tr3_nodes, "TD", { class: true });
      var td7_nodes = children(td7);
      a = claim_element(td7_nodes, "A", { href: true });
      var a_nodes = children(a);
      t18 = claim_text(
        a_nodes,
        /*site*/
        ctx[5]
      );
      a_nodes.forEach(detach);
      td7_nodes.forEach(detach);
      tr3_nodes.forEach(detach);
      t19 = claim_space(tbody_nodes);
      tr4 = claim_element(tbody_nodes, "TR", { class: true });
      var tr4_nodes = children(tr4);
      td8 = claim_element(tr4_nodes, "TD", { class: true, ["data-svelte-h"]: true });
      if (get_svelte_dataset(td8) !== "svelte-8hox4k")
        td8.textContent = textContent_4;
      t21 = claim_space(tr4_nodes);
      td9 = claim_element(tr4_nodes, "TD", { class: true });
      var td9_nodes = children(td9);
      t22 = claim_text(td9_nodes, t22_value);
      td9_nodes.forEach(detach);
      tr4_nodes.forEach(detach);
      tbody_nodes.forEach(detach);
      table_nodes.forEach(detach);
      this.h();
    },
    h() {
      attr(h2, "class", "svelte-19p5cuz");
      attr(div, "class", "name svelte-19p5cuz");
      attr(td0, "class", "svelte-19p5cuz");
      attr(td1, "class", "svelte-19p5cuz");
      attr(tr0, "class", "svelte-19p5cuz");
      attr(td2, "class", "svelte-19p5cuz");
      attr(td3, "class", "svelte-19p5cuz");
      attr(tr1, "class", "svelte-19p5cuz");
      attr(td4, "class", "svelte-19p5cuz");
      attr(td5, "class", "svelte-19p5cuz");
      attr(tr2, "class", "svelte-19p5cuz");
      attr(td6, "class", "svelte-19p5cuz");
      attr(
        a,
        "href",
        /*site*/
        ctx[5]
      );
      attr(td7, "class", "svelte-19p5cuz");
      attr(tr3, "class", "svelte-19p5cuz");
      attr(td8, "class", "svelte-19p5cuz");
      attr(td9, "class", "svelte-19p5cuz");
      attr(tr4, "class", "svelte-19p5cuz");
    },
    m(target, anchor) {
      insert_hydration(target, div, anchor);
      append_hydration(div, h2);
      append_hydration(h2, t0);
      append_hydration(div, t1);
      append_hydration(div, small);
      append_hydration(small, t2);
      insert_hydration(target, t3, anchor);
      insert_hydration(target, table, anchor);
      append_hydration(table, tbody);
      append_hydration(tbody, tr0);
      append_hydration(tr0, td0);
      append_hydration(tr0, t5);
      append_hydration(tr0, td1);
      append_hydration(td1, t6);
      append_hydration(tbody, t7);
      append_hydration(tbody, tr1);
      append_hydration(tr1, td2);
      append_hydration(tr1, t9);
      append_hydration(tr1, td3);
      append_hydration(td3, t10);
      append_hydration(tbody, t11);
      append_hydration(tbody, tr2);
      append_hydration(tr2, td4);
      append_hydration(tr2, t13);
      append_hydration(tr2, td5);
      append_hydration(td5, t14);
      append_hydration(tbody, t15);
      append_hydration(tbody, tr3);
      append_hydration(tr3, td6);
      append_hydration(tr3, t17);
      append_hydration(tr3, td7);
      append_hydration(td7, a);
      append_hydration(a, t18);
      append_hydration(tbody, t19);
      append_hydration(tbody, tr4);
      append_hydration(tr4, td8);
      append_hydration(tr4, t21);
      append_hydration(tr4, td9);
      append_hydration(td9, t22);
    },
    p(ctx2, [dirty]) {
      if (dirty & /*realname*/
      2)
        set_data(
          t0,
          /*realname*/
          ctx2[1]
        );
      if (dirty & /*occupation*/
      4)
        set_data(
          t2,
          /*occupation*/
          ctx2[2]
        );
      if (dirty & /*username*/
      1)
        set_data(
          t6,
          /*username*/
          ctx2[0]
        );
      if (dirty & /*age*/
      8)
        set_data(
          t10,
          /*age*/
          ctx2[3]
        );
      if (dirty & /*email*/
      16)
        set_data(
          t14,
          /*email*/
          ctx2[4]
        );
      if (dirty & /*site*/
      32)
        set_data(
          t18,
          /*site*/
          ctx2[5]
        );
      if (dirty & /*site*/
      32) {
        attr(
          a,
          "href",
          /*site*/
          ctx2[5]
        );
      }
      if (dirty & /*searching*/
      64 && t22_value !== (t22_value = /*searching*/
      ctx2[6] ? "Yes" : "No"))
        set_data(t22, t22_value);
    },
    i: noop,
    o: noop,
    d(detaching) {
      if (detaching) {
        detach(div);
        detach(t3);
        detach(table);
      }
    }
  };
}
function instance($$self, $$props, $$invalidate) {
  let { username } = $$props;
  let { realname } = $$props;
  let { occupation } = $$props;
  let { age } = $$props;
  let { email } = $$props;
  let { site } = $$props;
  let { searching } = $$props;
  $$self.$$set = ($$props2) => {
    if ("username" in $$props2)
      $$invalidate(0, username = $$props2.username);
    if ("realname" in $$props2)
      $$invalidate(1, realname = $$props2.realname);
    if ("occupation" in $$props2)
      $$invalidate(2, occupation = $$props2.occupation);
    if ("age" in $$props2)
      $$invalidate(3, age = $$props2.age);
    if ("email" in $$props2)
      $$invalidate(4, email = $$props2.email);
    if ("site" in $$props2)
      $$invalidate(5, site = $$props2.site);
    if ("searching" in $$props2)
      $$invalidate(6, searching = $$props2.searching);
  };
  return [username, realname, occupation, age, email, site, searching];
}
class Profile extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, instance, create_fragment, safe_not_equal, {
      username: 0,
      realname: 1,
      occupation: 2,
      age: 3,
      email: 4,
      site: 5,
      searching: 6
    });
  }
}
export {
  Profile as default
};
//# sourceMappingURL=profile-KyvgK15P.js.map
