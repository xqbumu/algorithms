import { A as AppState, g as golteContext, _ as __vitePreload } from "../chunks/appstate-tLgYl6N1.js";
import { S as SvelteComponent, i as init, s as safe_not_equal, c as construct_svelte_component, a as create_component, e as empty, b as claim_component, m as mount_component, d as insert_hydration, g as group_outros, t as transition_out, f as destroy_component, h as check_outros, j as transition_in, k as detach, l as component_subscribe, n as assign, o as noop, p as setContext, q as onMount, r as get_store_value } from "../chunks/index-6zOx8Fy8.js";
function get_spread_update(levels, updates) {
  const update = {};
  const to_null_out = {};
  const accounted_for = { $$scope: 1 };
  let i = levels.length;
  while (i--) {
    const o = levels[i];
    const n = updates[i];
    if (n) {
      for (const key in o) {
        if (!(key in n))
          to_null_out[key] = 1;
      }
      for (const key in n) {
        if (!accounted_for[key]) {
          update[key] = n[key];
          accounted_for[key] = 1;
        }
      }
      levels[i] = n;
    } else {
      for (const key in o) {
        accounted_for[key] = 1;
      }
    }
  }
  for (const key in to_null_out) {
    if (!(key in update))
      update[key] = void 0;
  }
  return update;
}
function get_spread_object(spread_props) {
  return typeof spread_props === "object" && spread_props !== null ? spread_props : {};
}
function create_if_block$1(ctx) {
  let node_1;
  let current;
  node_1 = new Node({
    props: {
      node: (
        /*$next*/
        ctx[1]
      ),
      index: (
        /*index*/
        ctx[0] + 1
      )
    }
  });
  return {
    c() {
      create_component(node_1.$$.fragment);
    },
    l(nodes) {
      claim_component(node_1.$$.fragment, nodes);
    },
    m(target, anchor) {
      mount_component(node_1, target, anchor);
      current = true;
    },
    p(ctx2, dirty) {
      const node_1_changes = {};
      if (dirty & /*$next*/
      2)
        node_1_changes.node = /*$next*/
        ctx2[1];
      if (dirty & /*index*/
      1)
        node_1_changes.index = /*index*/
        ctx2[0] + 1;
      node_1.$set(node_1_changes);
    },
    i(local) {
      if (current)
        return;
      transition_in(node_1.$$.fragment, local);
      current = true;
    },
    o(local) {
      transition_out(node_1.$$.fragment, local);
      current = false;
    },
    d(detaching) {
      destroy_component(node_1, detaching);
    }
  };
}
function create_key_block$1(ctx) {
  let if_block_anchor;
  let current;
  let if_block = (
    /*$next*/
    ctx[1] && create_if_block$1(ctx)
  );
  return {
    c() {
      if (if_block)
        if_block.c();
      if_block_anchor = empty();
    },
    l(nodes) {
      if (if_block)
        if_block.l(nodes);
      if_block_anchor = empty();
    },
    m(target, anchor) {
      if (if_block)
        if_block.m(target, anchor);
      insert_hydration(target, if_block_anchor, anchor);
      current = true;
    },
    p(ctx2, dirty) {
      if (
        /*$next*/
        ctx2[1]
      ) {
        if (if_block) {
          if_block.p(ctx2, dirty);
          if (dirty & /*$next*/
          2) {
            transition_in(if_block, 1);
          }
        } else {
          if_block = create_if_block$1(ctx2);
          if_block.c();
          transition_in(if_block, 1);
          if_block.m(if_block_anchor.parentNode, if_block_anchor);
        }
      } else if (if_block) {
        group_outros();
        transition_out(if_block, 1, 1, () => {
          if_block = null;
        });
        check_outros();
      }
    },
    i(local) {
      if (current)
        return;
      transition_in(if_block);
      current = true;
    },
    o(local) {
      transition_out(if_block);
      current = false;
    },
    d(detaching) {
      if (detaching) {
        detach(if_block_anchor);
      }
      if (if_block)
        if_block.d(detaching);
    }
  };
}
function create_default_slot(ctx) {
  let previous_key = (
    /*$next*/
    ctx[1]
  );
  let key_block_anchor;
  let current;
  let key_block = create_key_block$1(ctx);
  return {
    c() {
      key_block.c();
      key_block_anchor = empty();
    },
    l(nodes) {
      key_block.l(nodes);
      key_block_anchor = empty();
    },
    m(target, anchor) {
      key_block.m(target, anchor);
      insert_hydration(target, key_block_anchor, anchor);
      current = true;
    },
    p(ctx2, dirty) {
      if (dirty & /*$next*/
      2 && safe_not_equal(previous_key, previous_key = /*$next*/
      ctx2[1])) {
        group_outros();
        transition_out(key_block, 1, 1, noop);
        check_outros();
        key_block = create_key_block$1(ctx2);
        key_block.c();
        transition_in(key_block, 1);
        key_block.m(key_block_anchor.parentNode, key_block_anchor);
      } else {
        key_block.p(ctx2, dirty);
      }
    },
    i(local) {
      if (current)
        return;
      transition_in(key_block);
      current = true;
    },
    o(local) {
      transition_out(key_block);
      current = false;
    },
    d(detaching) {
      if (detaching) {
        detach(key_block_anchor);
      }
      key_block.d(detaching);
    }
  };
}
function create_fragment$1(ctx) {
  let switch_instance;
  let switch_instance_anchor;
  let current;
  const switch_instance_spread_levels = [
    /*content*/
    ctx[3].props
  ];
  var switch_value = (
    /*content*/
    ctx[3].comp
  );
  function switch_props(ctx2, dirty) {
    let switch_instance_props = {
      $$slots: { default: [create_default_slot] },
      $$scope: { ctx: ctx2 }
    };
    if (dirty !== void 0 && dirty & /*content*/
    8) {
      switch_instance_props = get_spread_update(switch_instance_spread_levels, [get_spread_object(
        /*content*/
        ctx2[3].props
      )]);
    } else {
      for (let i = 0; i < switch_instance_spread_levels.length; i += 1) {
        switch_instance_props = assign(switch_instance_props, switch_instance_spread_levels[i]);
      }
    }
    return { props: switch_instance_props };
  }
  if (switch_value) {
    switch_instance = construct_svelte_component(switch_value, switch_props(ctx));
  }
  return {
    c() {
      if (switch_instance)
        create_component(switch_instance.$$.fragment);
      switch_instance_anchor = empty();
    },
    l(nodes) {
      if (switch_instance)
        claim_component(switch_instance.$$.fragment, nodes);
      switch_instance_anchor = empty();
    },
    m(target, anchor) {
      if (switch_instance)
        mount_component(switch_instance, target, anchor);
      insert_hydration(target, switch_instance_anchor, anchor);
      current = true;
    },
    p(ctx2, [dirty]) {
      if (switch_value !== (switch_value = /*content*/
      ctx2[3].comp)) {
        if (switch_instance) {
          group_outros();
          const old_component = switch_instance;
          transition_out(old_component.$$.fragment, 1, 0, () => {
            destroy_component(old_component, 1);
          });
          check_outros();
        }
        if (switch_value) {
          switch_instance = construct_svelte_component(switch_value, switch_props(ctx2, dirty));
          create_component(switch_instance.$$.fragment);
          transition_in(switch_instance.$$.fragment, 1);
          mount_component(switch_instance, switch_instance_anchor.parentNode, switch_instance_anchor);
        } else {
          switch_instance = null;
        }
      } else if (switch_value) {
        const switch_instance_changes = dirty & /*content*/
        8 ? get_spread_update(switch_instance_spread_levels, [get_spread_object(
          /*content*/
          ctx2[3].props
        )]) : {};
        if (dirty & /*$$scope, $next, index*/
        35) {
          switch_instance_changes.$$scope = { dirty, ctx: ctx2 };
        }
        switch_instance.$set(switch_instance_changes);
      }
    },
    i(local) {
      if (current)
        return;
      if (switch_instance)
        transition_in(switch_instance.$$.fragment, local);
      current = true;
    },
    o(local) {
      if (switch_instance)
        transition_out(switch_instance.$$.fragment, local);
      current = false;
    },
    d(detaching) {
      if (detaching) {
        detach(switch_instance_anchor);
      }
      if (switch_instance)
        destroy_component(switch_instance, detaching);
    }
  };
}
function instance$1($$self, $$props, $$invalidate) {
  let $next;
  let { node } = $$props;
  let { index } = $$props;
  const { next, content } = node;
  component_subscribe($$self, next, (value) => $$invalidate(1, $next = value));
  $$self.$$set = ($$props2) => {
    if ("node" in $$props2)
      $$invalidate(4, node = $$props2.node);
    if ("index" in $$props2)
      $$invalidate(0, index = $$props2.index);
  };
  return [index, $next, next, content, node];
}
class Node_1 extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, instance$1, create_fragment$1, safe_not_equal, { node: 4, index: 0 });
  }
}
function csrWrapper(options) {
  const ssrError = options.props.node.content.ssrError;
  if (ssrError)
    return new options.props.node.content.errPage({ ...options, props: ssrError });
  return new Node_1(options);
}
const Node = csrWrapper;
function create_if_block(ctx) {
  let node_1;
  let current;
  node_1 = new Node({
    props: { node: (
      /*$node*/
      ctx[0]
    ), index: 0 }
  });
  return {
    c() {
      create_component(node_1.$$.fragment);
    },
    l(nodes) {
      claim_component(node_1.$$.fragment, nodes);
    },
    m(target, anchor) {
      mount_component(node_1, target, anchor);
      current = true;
    },
    p(ctx2, dirty) {
      const node_1_changes = {};
      if (dirty & /*$node*/
      1)
        node_1_changes.node = /*$node*/
        ctx2[0];
      node_1.$set(node_1_changes);
    },
    i(local) {
      if (current)
        return;
      transition_in(node_1.$$.fragment, local);
      current = true;
    },
    o(local) {
      transition_out(node_1.$$.fragment, local);
      current = false;
    },
    d(detaching) {
      destroy_component(node_1, detaching);
    }
  };
}
function create_key_block(ctx) {
  let if_block_anchor;
  let current;
  let if_block = (
    /*$node*/
    ctx[0] && create_if_block(ctx)
  );
  return {
    c() {
      if (if_block)
        if_block.c();
      if_block_anchor = empty();
    },
    l(nodes) {
      if (if_block)
        if_block.l(nodes);
      if_block_anchor = empty();
    },
    m(target, anchor) {
      if (if_block)
        if_block.m(target, anchor);
      insert_hydration(target, if_block_anchor, anchor);
      current = true;
    },
    p(ctx2, dirty) {
      if (
        /*$node*/
        ctx2[0]
      ) {
        if (if_block) {
          if_block.p(ctx2, dirty);
          if (dirty & /*$node*/
          1) {
            transition_in(if_block, 1);
          }
        } else {
          if_block = create_if_block(ctx2);
          if_block.c();
          transition_in(if_block, 1);
          if_block.m(if_block_anchor.parentNode, if_block_anchor);
        }
      } else if (if_block) {
        group_outros();
        transition_out(if_block, 1, 1, () => {
          if_block = null;
        });
        check_outros();
      }
    },
    i(local) {
      if (current)
        return;
      transition_in(if_block);
      current = true;
    },
    o(local) {
      transition_out(if_block);
      current = false;
    },
    d(detaching) {
      if (detaching) {
        detach(if_block_anchor);
      }
      if (if_block)
        if_block.d(detaching);
    }
  };
}
function create_fragment(ctx) {
  let previous_key = (
    /*$node*/
    ctx[0]
  );
  let key_block_anchor;
  let current;
  let key_block = create_key_block(ctx);
  return {
    c() {
      key_block.c();
      key_block_anchor = empty();
    },
    l(nodes) {
      key_block.l(nodes);
      key_block_anchor = empty();
    },
    m(target, anchor) {
      key_block.m(target, anchor);
      insert_hydration(target, key_block_anchor, anchor);
      current = true;
    },
    p(ctx2, [dirty]) {
      if (dirty & /*$node*/
      1 && safe_not_equal(previous_key, previous_key = /*$node*/
      ctx2[0])) {
        group_outros();
        transition_out(key_block, 1, 1, noop);
        check_outros();
        key_block = create_key_block(ctx2);
        key_block.c();
        transition_in(key_block, 1);
        key_block.m(key_block_anchor.parentNode, key_block_anchor);
      } else {
        key_block.p(ctx2, dirty);
      }
    },
    i(local) {
      if (current)
        return;
      transition_in(key_block);
      current = true;
    },
    o(local) {
      transition_out(key_block);
      current = false;
    },
    d(detaching) {
      if (detaching) {
        detach(key_block_anchor);
      }
      key_block.d(detaching);
    }
  };
}
function instance($$self, $$props, $$invalidate) {
  let $node;
  let { nodes } = $$props;
  let { contextData } = $$props;
  const state = new AppState(contextData.URL, nodes);
  const { node } = state;
  component_subscribe($$self, node, (value) => $$invalidate(0, $node = value));
  setContext(golteContext, state);
  onMount(() => {
    history.replaceState(get_store_value(state.url).href, "");
    addEventListener("popstate", async (e) => {
      if (!e.state)
        return;
      await state.update(e.state);
    });
  });
  $$self.$$set = ($$props2) => {
    if ("nodes" in $$props2)
      $$invalidate(2, nodes = $$props2.nodes);
    if ("contextData" in $$props2)
      $$invalidate(3, contextData = $$props2.contextData);
  };
  return [$node, node, nodes, contextData];
}
class Root extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, instance, create_fragment, safe_not_equal, { nodes: 2, contextData: 3 });
  }
}
async function hydrate(target, nodes, contextData) {
  const promise = Promise.all(nodes.map(async (n) => ({
    comp: (await __vitePreload(() => import(n.comp), true ? __vite__mapDeps([]) : void 0)).default,
    props: n.props,
    errPage: (await __vitePreload(() => import(n.errPage), true ? __vite__mapDeps([]) : void 0)).default,
    ssrError: n.ssrError
  })));
  new Root({
    target,
    props: {
      nodes: await promise,
      contextData
    },
    hydrate: true
  });
}
export {
  hydrate
};
function __vite__mapDeps(indexes) {
  if (!__vite__mapDeps.viteFileDeps) {
    __vite__mapDeps.viteFileDeps = []
  }
  return indexes.map((i) => __vite__mapDeps.viteFileDeps[i])
}
//# sourceMappingURL=hydrate--REBL6VH.js.map
