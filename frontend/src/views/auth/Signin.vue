<template>
  <div class="mx-auto w-full max-w-sm">
    <img
      class="h-12 w-auto mx-auto mb-8"
      src="../../assets/logo-full.svg"
      alt="Bytebase"
    />

    <div class="mt-8">
      <n-card>
        <n-tabs
          class="card-tabs"
          default-value="standard"
          size="small"
          animated
          pane-style="padding: 12px 0 0 0"
        >
          <n-tab-pane name="standard" tab="Standard">
            <form class="space-y-6 px-1" @submit.prevent="trySignin()">
              <div>
                <label
                  for="email"
                  class="block text-sm font-medium leading-5 text-control"
                >
                  {{ $t("common.email") }}
                  <span class="text-red-600">*</span>
                </label>
                <div class="mt-1 rounded-md shadow-sm">
                  <BBTextField
                    v-model:value="state.email"
                    required
                    :input-props="{
                      id: 'email',
                      autocomplete: 'on',
                      type: 'email',
                    }"
                    placeholder="jim@example.com"
                  />
                </div>
              </div>

              <div>
                <label
                  for="password"
                  class="flex justify-between text-sm font-medium leading-5 text-control"
                >
                  <div>
                    {{ $t("common.password") }}
                    <span class="text-red-600">*</span>
                  </div>
                  <router-link
                    :to="{
                      path: '/auth/password-forgot',
                      query: {
                        hint: route.query.hint,
                      },
                    }"
                    class="text-sm font-normal text-control-light hover:underline focus:outline-none"
                    tabindex="-1"
                  >
                    {{ $t("auth.sign-in.forget-password") }}
                  </router-link>
                </label>
                <div
                  class="relative flex flex-row items-center mt-1 rounded-md shadow-sm"
                >
                  <BBTextField
                    v-model:value="state.password"
                    :type="state.showPassword ? 'text' : 'password'"
                    :input-props="{ id: 'password', autocomplete: 'on' }"
                    required
                  />
                  <div
                    class="hover:cursor-pointer absolute right-3"
                    @click="
                      () => {
                        state.showPassword = !state.showPassword;
                      }
                    "
                  >
                    <heroicons-outline:eye
                      v-if="state.showPassword"
                      class="w-4 h-4"
                    />
                    <heroicons-outline:eye-slash v-else class="w-4 h-4" />
                  </div>
                </div>
              </div>

              <div class="w-full">
                <NButton
                  attr-type="submit"
                  type="primary"
                  :disabled="!allowSignin()"
                  size="large"
                  style="width: 100%"
                >
                  {{ $t("common.sign-in") }}
                </NButton>
              </div>
            </form>

            <div class="mt-3">
              <div
                class="flex justify-center items-center text-sm text-control"
              >
                <template v-if="isDemo">
                  <span class="text-accent">
                    {{
                      $t("auth.sign-in.demo-note", {
                        username: "demo@example.com",
                        password: "1024",
                      })
                    }}
                  </span>
                </template>
                <template v-else-if="!disallowSignup">
                  <span>
                    {{ $t("auth.sign-in.new-user") }}
                  </span>
                  <router-link to="/auth/signup" class="accent-link px-2">
                    {{ $t("common.sign-up") }}
                  </router-link>
                </template>
              </div>
            </div>
          </n-tab-pane>

          <template
            v-for="identityProvider in groupedIdentityProviderList"
            :key="identityProvider.name"
          >
            <n-tab-pane
              v-if="identityProvider.type === IdentityProviderType.LDAP"
              :name="identityProvider.name"
              :tab="identityProvider.title"
            >
              <form
                class="space-y-6"
                @submit.prevent="trySignin(identityProvider.name)"
              >
                <div>
                  <label
                    for="email"
                    class="block text-sm font-medium leading-5 text-control"
                  >
                    {{ $t("common.username") }}
                    <span class="text-red-600">*</span>
                  </label>
                  <div class="mt-1 rounded-md shadow-sm">
                    <BBTextField
                      v-model:value="state.email"
                      required
                      placeholder="jim"
                      :input-props="{ id: 'username', autocomplete: 'on' }"
                    />
                  </div>
                </div>

                <div>
                  <label
                    for="password"
                    class="flex justify-between text-sm font-medium leading-5 text-control"
                  >
                    <div>
                      {{ $t("common.password") }}
                      <span class="text-red-600">*</span>
                    </div>
                  </label>
                  <div
                    class="relative flex flex-row items-center mt-1 rounded-md shadow-sm"
                  >
                    <BBTextField
                      v-model:value="state.password"
                      :type="state.showPassword ? 'text' : 'password'"
                      :input-props="{ id: 'password', autocomplete: 'on' }"
                      required
                    />
                    <div
                      class="hover:cursor-pointer absolute right-3"
                      @click="
                        () => {
                          state.showPassword = !state.showPassword;
                        }
                      "
                    >
                      <heroicons-outline:eye
                        v-if="state.showPassword"
                        class="w-4 h-4"
                      />
                      <heroicons-outline:eye-slash v-else class="w-4 h-4" />
                    </div>
                  </div>
                </div>

                <div class="w-full">
                  <NButton
                    attr-type="submit"
                    type="primary"
                    size="large"
                    :disabled="!allowSignin(identityProvider.name)"
                    style="width: 100%"
                  >
                    {{ $t("common.sign-in") }}
                  </NButton>
                </div>
              </form>
            </n-tab-pane>
          </template>
        </n-tabs>
      </n-card>
    </div>

    <div v-if="separatedIdentityProviderList.length > 0" class="mb-3 px-1">
      <div class="relative my-4">
        <div class="absolute inset-0 flex items-center" aria-hidden="true">
          <div class="w-full border-t border-control-border"></div>
        </div>
        <div class="relative flex justify-center text-sm">
          <span class="px-2 bg-white text-control">{{ $t("common.or") }}</span>
        </div>
      </div>
      <template
        v-for="identityProvider in separatedIdentityProviderList"
        :key="identityProvider.name"
      >
        <div class="w-full mb-2">
          <NButton
            style="width: 100%"
            size="large"
            @click.prevent="trySigninWithIdentityProvider(identityProvider)"
          >
            {{
              $t("auth.sign-in.sign-in-with-idp", {
                idp: identityProvider.title,
              })
            }}
          </NButton>
        </div>
      </template>
    </div>
  </div>
  <AuthFooter />
</template>

<script lang="ts" setup>
import { storeToRefs } from "pinia";
import { computed, onMounted, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import { AUTH_MFA_MODULE, AUTH_SIGNUP_MODULE } from "@/router/auth";
import {
  useActuatorV1Store,
  useAuthStore,
  useIdentityProviderStore,
} from "@/store";
import {
  IdentityProvider,
  IdentityProviderType,
} from "@/types/proto/v1/idp_service";
import { isValidEmail, openWindowForSSO } from "@/utils";
import AuthFooter from "./AuthFooter.vue";

interface LocalState {
  email: string;
  password: string;
  showPassword: boolean;
}

const router = useRouter();
const route = useRoute();
const actuatorStore = useActuatorV1Store();
const authStore = useAuthStore();
const identityProviderStore = useIdentityProviderStore();

const state = reactive<LocalState>({
  email: "",
  password: "",
  showPassword: false,
});
const { isDemo, disallowSignup } = storeToRefs(actuatorStore);

const separatedIdentityProviderList = computed(() =>
  identityProviderStore.identityProviderList.filter(
    (idp) => idp.type !== IdentityProviderType.LDAP
  )
);
const groupedIdentityProviderList = computed(() =>
  identityProviderStore.identityProviderList.filter(
    (idp) => idp.type === IdentityProviderType.LDAP
  )
);

onMounted(async () => {
  // Navigate to signup if needs admin setup.
  // Unable to achieve it in router.beforeEach because actuator/info is fetched async and returns
  // after router has already made the decision on first page load.
  if (actuatorStore.needAdminSetup && !disallowSignup.value) {
    router.push({ name: AUTH_SIGNUP_MODULE, replace: true });
  }

  const url = new URL(window.location.href);
  const params = new URLSearchParams(url.search);
  state.email = params.get("email") ?? (isDemo.value ? "demo@example.com" : "");
  state.password = params.get("password") ?? (isDemo.value ? "1024" : "");
  state.showPassword = !!isDemo.value;

  await identityProviderStore.fetchIdentityProviderList();
  if (
    window.location.href.startsWith("https://demo.bytebase.com") &&
    isDemo.value &&
    state.email &&
    state.password
  ) {
    await trySignin();
  }
});

const allowSignin = (idpName?: string) => {
  if (!idpName) {
    return isValidEmail(state.email) && state.password;
  }
  return state.email && state.password;
};

const trySignin = async (idpName?: string) => {
  const mfaTempToken = await authStore.login({
    email: state.email,
    password: state.password,
    web: true,
    idpName: idpName,
  });
  if (mfaTempToken) {
    router.push({
      name: AUTH_MFA_MODULE,
      query: {
        mfaTempToken,
        redirect: route.query.redirect as string,
      },
    });
  } else {
    router.push("/");
  }
};

const trySigninWithIdentityProvider = async (
  identityProvider: IdentityProvider
) => {
  await openWindowForSSO(
    identityProvider,
    false /* !popup */,
    route.query.redirect as string
  );
};
</script>
