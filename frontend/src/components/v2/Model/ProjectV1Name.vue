<template>
  <component
    :is="link ? 'router-link' : tag"
    v-bind="bindings"
    class="inline-flex items-center gap-x-1"
    :class="link && !plain && 'normal-link'"
  >
    <NPerformantEllipsis :line-clamp="1">
      {{ projectV1Name(project) }}
    </NPerformantEllipsis>
  </component>
</template>

<script lang="ts" setup>
import { NPerformantEllipsis } from "naive-ui";
import { computed } from "vue";
import type { Project } from "@/types/proto/v1/project_service";
import { projectV1Name } from "@/utils";

const props = withDefaults(
  defineProps<{
    project: Project;
    tag?: string;
    link?: boolean;
    plain?: boolean;
    hash?: string;
  }>(),
  {
    tag: "span",
    link: true,
    plain: false,
    hash: "",
  }
);

const bindings = computed(() => {
  if (props.link) {
    return {
      to: `/${props.project.name}`,
      activeClass: "",
      exactActiveClass: "",
      onClick: (e: MouseEvent) => {
        e.stopPropagation();
      },
    };
  }
  return {};
});
</script>
