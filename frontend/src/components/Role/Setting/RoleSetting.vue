<template>
  <div class="space-y-4">
    <div class="w-full flex flex-row justify-between items-center">
      <div class="textinfolabel">
        {{ $t("role.setting.description") }}
        <a
          href="https://www.bytebase.com/docs/administration/custom-roles?source=console"
          class="normal-link text-sm inline-flex flex-row items-center"
          target="_blank"
        >
          {{ $t("common.learn-more") }}
          <heroicons-outline:external-link class="w-4 h-4" />
        </a>
      </div>
      <NButton type="primary" :disabled="!allowCreateRole" @click="addRole">
        {{ $t("role.setting.add") }}
      </NButton>
    </div>
    <RoleDataTable
      :role-list="filteredRoleList"
      :show-placeholder="state.ready"
      @select-role="selectRole($event, 'EDIT')"
    />

    <div
      v-if="!state.ready"
      class="relative flex flex-col h-[8rem] items-center justify-center"
    >
      <BBSpin />
    </div>

    <RolePanel
      :role="state.detail.role"
      :mode="state.detail.mode"
      @close="state.detail.role = undefined"
    />

    <FeatureModal
      feature="bb.feature.custom-role"
      :open="showFeatureModal"
      @cancel="showFeatureModal = false"
    />
  </div>
</template>

<script lang="ts" setup>
import { NButton } from "naive-ui";
import { computed, onMounted, reactive, ref } from "vue";
import { featureToRef, useCurrentUserV1, useRoleStore } from "@/store";
import { Role } from "@/types/proto/v1/role_service";
import { hasWorkspacePermissionV2 } from "@/utils";
import { RoleDataTable, RolePanel } from "./components";
import { provideCustomRoleSettingContext } from "./context";

type LocalState = {
  ready: boolean;
  detail: {
    role: Role | undefined;
    mode: "ADD" | "EDIT";
  };
  filter: {
    keyword: string;
  };
};

const roleStore = useRoleStore();
const state = reactive<LocalState>({
  ready: false,
  detail: {
    role: undefined,
    mode: "ADD",
  },
  filter: {
    keyword: "",
  },
});

const currentUser = useCurrentUserV1();
const hasCustomRoleFeature = featureToRef("bb.feature.custom-role");
const showFeatureModal = ref(false);

const allowCreateRole = computed(() => {
  return hasWorkspacePermissionV2(currentUser.value, "bb.roles.create");
});

const filteredRoleList = computed(() => {
  const keyword = state.filter.keyword.trim().toLowerCase();
  if (!keyword) return roleStore.roleList;
  return roleStore.roleList.filter((role) => {
    return (
      role.name.toLowerCase().includes(keyword) ||
      role.description.toLowerCase().includes(keyword)
    );
  });
});

const addRole = () => {
  selectRole(Role.fromJSON({}), "ADD");
};

const selectRole = (role: Role | undefined, mode?: "ADD" | "EDIT") => {
  state.detail.role = role;
  if (mode) {
    state.detail.mode = mode;
  }
};

const prepare = async () => {
  try {
    await roleStore.fetchRoleList();
  } finally {
    state.ready = true;
  }
};
onMounted(prepare);

provideCustomRoleSettingContext({
  hasCustomRoleFeature,
  showFeatureModal,
});
</script>
