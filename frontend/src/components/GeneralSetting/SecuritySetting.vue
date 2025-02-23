<template>
  <div class="py-6 lg:flex">
    <div class="text-left lg:w-1/4">
      <div class="flex items-center space-x-2">
        <h1 class="text-2xl font-bold">
          {{ $t("settings.general.workspace.security") }}
        </h1>
        <FeatureBadge feature="bb.feature.watermark" />
      </div>
      <span v-if="!allowEdit" class="text-sm text-gray-400">
        {{ $t("settings.general.workspace.only-admin-can-edit") }}
      </span>
    </div>
    <div class="flex-1 lg:px-4">
      <div class="mb-7 mt-4 lg:mt-0">
        <NTooltip placement="top-start" :disabled="allowEdit">
          <template #trigger>
            <label
              class="flex items-center gap-x-2"
              :class="[allowEdit ? 'cursor-pointer' : 'cursor-not-allowed']"
            >
              <NCheckbox
                :disabled="!allowEdit"
                :checked="watermarkEnabled"
                :label="$t('settings.general.workspace.watermark.enable')"
                @update:checked="handleWatermarkToggle"
              />
            </label>
          </template>
          <span class="text-sm text-gray-400 -translate-y-2">
            {{ $t("settings.general.workspace.watermark.only-admin-can-edit") }}
          </span>
        </NTooltip>
        <div class="mb-3 text-sm text-gray-400">
          {{ $t("settings.general.workspace.watermark.description") }}
        </div>
      </div>
      <div v-if="!isSaaSMode" class="mb-7 mt-4 lg:mt-0">
        <NTooltip placement="top-start" :disabled="allowEdit">
          <template #trigger>
            <label
              class="flex items-center gap-x-2"
              :class="[allowEdit ? 'cursor-pointer' : 'cursor-not-allowed']"
            >
              <NCheckbox
                :disabled="!allowEdit"
                :checked="disallowSignupEnabled"
                :label="$t('settings.general.workspace.disallow-signup.enable')"
                @update:checked="handleDisallowSignupToggle"
              />
            </label>
          </template>
          <span class="text-sm text-gray-400 -translate-y-2">
            {{ $t("settings.general.workspace.only-admin-can-edit") }}
          </span>
        </NTooltip>
        <div class="mb-3 text-sm text-gray-400">
          {{ $t("settings.general.workspace.disallow-signup.description") }}
        </div>
      </div>
      <div class="mb-7 mt-4 lg:mt-0">
        <NTooltip placement="top-start" :disabled="allowEdit">
          <template #trigger>
            <label
              class="flex items-center gap-x-2"
              :class="[allowEdit ? 'cursor-pointer' : 'cursor-not-allowed']"
            >
              <NCheckbox
                :disabled="!allowEdit"
                :checked="require2FAEnabled"
                :label="$t('settings.general.workspace.require-2fa.enable')"
                @update:checked="handleRequire2FAToggle"
              />
            </label>
          </template>
          <span class="text-sm text-gray-400 -translate-y-2">
            {{ $t("settings.general.workspace.only-admin-can-edit") }}
          </span>
        </NTooltip>
        <div class="mb-3 text-sm text-gray-400">
          {{ $t("settings.general.workspace.require-2fa.description") }}
        </div>
      </div>
      <div class="mb-7 mt-4 lg:mt-0">
        <NTooltip placement="top-start" :disabled="allowEdit">
          <template #trigger>
            <label
              class="flex items-center gap-x-2"
              :class="[allowEdit ? 'cursor-pointer' : 'cursor-not-allowed']"
            >
              <NCheckbox
                :disabled="!allowEdit"
                :checked="restrictIssueCreationForSQLReview"
                :label="
                  $t(
                    'settings.general.workspace.restrict-issue-creation-for-sql-review.title'
                  )
                "
                @update:checked="handleRestrictIssueCreationForSQLReviewToggle"
              />
            </label>
          </template>
          <span class="text-sm text-gray-400 -translate-y-2">
            {{ $t("settings.general.workspace.only-admin-can-edit") }}
          </span>
        </NTooltip>
        <div class="mb-3 text-sm text-gray-400">
          {{
            $t(
              "settings.general.workspace.restrict-issue-creation-for-sql-review.description"
            )
          }}
        </div>
      </div>
      <SignInFrequencySetting :allow-edit="allowEdit" />
    </div>
  </div>

  <FeatureModal
    :open="!!state.featureNameForModal"
    :feature="state.featureNameForModal"
    @cancel="state.featureNameForModal = undefined"
  />
</template>

<script lang="ts" setup>
import { NCheckbox } from "naive-ui";
import { storeToRefs } from "pinia";
import { computed, onMounted, reactive } from "vue";
import { useI18n } from "vue-i18n";
import {
  featureToRef,
  pushNotification,
  useActuatorV1Store,
  usePolicyV1Store,
} from "@/store";
import { useSettingV1Store } from "@/store/modules/v1/setting";
import { FeatureType } from "@/types";
import {
  PolicyResourceType,
  PolicyType,
} from "@/types/proto/v1/org_policy_service";
import SignInFrequencySetting from "./SignInFrequencySetting.vue";

interface LocalState {
  featureNameForModal?: FeatureType;
}

defineProps<{
  allowEdit: boolean;
}>();

const state = reactive<LocalState>({});
const { t } = useI18n();
const settingV1Store = useSettingV1Store();
const actuatorStore = useActuatorV1Store();
const policyV1Store = usePolicyV1Store();

const { isSaaSMode } = storeToRefs(actuatorStore);
const hasWatermarkFeature = featureToRef("bb.feature.branding");
const has2FAFeature = featureToRef("bb.feature.2fa");
const hasDisallowSignupFeature = featureToRef("bb.feature.disallow-signup");

const watermarkEnabled = computed((): boolean => {
  return (
    settingV1Store.getSettingByName("bb.workspace.watermark")?.value
      ?.stringValue === "1"
  );
});
const disallowSignupEnabled = computed((): boolean => {
  return settingV1Store.workspaceProfileSetting?.disallowSignup ?? false;
});
const require2FAEnabled = computed((): boolean => {
  return settingV1Store.workspaceProfileSetting?.require2fa ?? false;
});
const restrictIssueCreationForSQLReview = computed((): boolean => {
  return (
    policyV1Store.getPolicyByName(
      "policies/RESTRICT_ISSUE_CREATION_FOR_SQL_REVIEW"
    )?.restrictIssueCreationForSqlReviewPolicy?.disallow ?? false
  );
});

onMounted(async () => {
  await prepareOrgPolicy();
});

const prepareOrgPolicy = async () => {
  await policyV1Store.getOrFetchPolicyByName(
    "policies/RESTRICT_ISSUE_CREATION_FOR_SQL_REVIEW"
  );
};

const handleDisallowSignupToggle = async (on: boolean) => {
  if (!hasDisallowSignupFeature.value && on) {
    state.featureNameForModal = "bb.feature.disallow-signup";
    return;
  }
  await settingV1Store.updateWorkspaceProfile({
    disallowSignup: on,
  });
  pushNotification({
    module: "bytebase",
    style: "SUCCESS",
    title: t("settings.general.workspace.config-updated"),
  });
};

const handleRequire2FAToggle = async (on: boolean) => {
  if (!has2FAFeature.value && on) {
    state.featureNameForModal = "bb.feature.2fa";
    return;
  }

  await settingV1Store.updateWorkspaceProfile({
    require2fa: on,
  });
  pushNotification({
    module: "bytebase",
    style: "SUCCESS",
    title: t("settings.general.workspace.config-updated"),
  });
};

const handleWatermarkToggle = async (on: boolean) => {
  if (!hasWatermarkFeature.value && on) {
    state.featureNameForModal = "bb.feature.watermark";
    return;
  }
  const value = on ? "1" : "0";
  await settingV1Store.upsertSetting({
    name: "bb.workspace.watermark",
    value: {
      stringValue: value,
    },
  });
  pushNotification({
    module: "bytebase",
    style: "SUCCESS",
    title: t("settings.general.workspace.watermark.update-success"),
  });
};

const handleRestrictIssueCreationForSQLReviewToggle = async (on: boolean) => {
  await policyV1Store.createPolicy("", {
    type: PolicyType.RESTRICT_ISSUE_CREATION_FOR_SQL_REVIEW,
    resourceType: PolicyResourceType.WORKSPACE,
    restrictIssueCreationForSqlReviewPolicy: {
      disallow: on,
    },
  });

  pushNotification({
    module: "bytebase",
    style: "SUCCESS",
    title: t("settings.general.workspace.config-updated"),
  });
};
</script>
