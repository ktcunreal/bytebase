<template>
  <BBGrid
    :column-list="columnList"
    :data-source="projectList"
    class="border"
    :show-placeholder="true"
    @click-row="clickProject"
  >
    <template #item="{ item: project }: ProjectGridRow">
      <div v-if="currentProject" class="bb-grid-cell">
        <CheckIcon
          v-if="currentProject.name === project.name"
          class="w-4 text-accent"
        />
      </div>
      <div class="bb-grid-cell text-gray-500">
        <span class="flex flex-row items-center space-x-1">
          <span>{{ project.key }}</span>
        </span>
      </div>
      <div class="bb-grid-cell truncate">
        <ProjectCol
          mode="ALL_SHORT"
          :project="project"
          :show-tenant-icon="true"
        />
      </div>
    </template>
  </BBGrid>
</template>

<script lang="ts" setup>
import { CheckIcon } from "lucide-vue-next";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { BBGridColumn, BBGridRow, BBGrid } from "@/bbkit";
import {
  PROJECT_V1_ROUTE_DETAIL,
  PROJECT_V1_ROUTE,
} from "@/router/dashboard/projectV1";
import { getProjectName } from "@/store/modules/v1/common";
import { Project } from "@/types/proto/v1/project_service";

export type ProjectGridRow = BBGridRow<Project>;

const props = defineProps<{
  projectList: Project[];
  currentProject?: Project;
}>();

const emit = defineEmits<{
  (event: "click"): void;
}>();

const router = useRouter();
const { t } = useI18n();
const columnList = computed((): BBGridColumn[] => {
  const list = [
    {
      title: t("project.table.key"),
      width: "minmax(min-content, 25%)",
    },
    {
      title: t("project.table.name"),
      width: "minmax(min-content, auto)",
    },
  ];
  if (props.currentProject) {
    list.unshift({
      title: "",
      width: "minmax(min-content, 3rem)",
    });
  }
  return list;
});

const clickProject = function (
  project: Project,
  section: number,
  row: number,
  e: MouseEvent
) {
  let routeName = PROJECT_V1_ROUTE_DETAIL;
  if (router.currentRoute.value.name?.toString().startsWith(PROJECT_V1_ROUTE)) {
    routeName = router.currentRoute.value.name?.toString();
  }

  const route = router.resolve({
    name: routeName,
    params: {
      projectId: getProjectName(project.name),
    },
  });

  if (e.ctrlKey || e.metaKey) {
    window.open(route.fullPath, "_blank");
  } else {
    window.location.href = route.fullPath;
  }
  emit("click");
};
</script>
