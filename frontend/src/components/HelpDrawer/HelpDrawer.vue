<template>
  <Drawer
    :show="active && shouldShowHelpDrawer"
    class="!w-96 max-w-full"
    @update:show="(show: boolean) => !show && onClose()"
  >
    <DrawerContent
      class="w-full"
      :title="state.frontmatter.title"
      :closable="true"
    >
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-if="state.html" class="overflow-auto" v-html="state.html"></div>
      <template #footer>
        <div class="flex flex-row justify-center pb-10">
          <div
            v-if="locale === 'zh-CN'"
            class="w-full flex flex-col items-center pt-2"
          >
            <p class="text-sm mb-2">微信扫码加入官方社群</p>
            <div class="flex flex-row justify-center">
              <div class="qrcode-card">
                <img
                  src="@/assets/wechat-official-qrcode.webp"
                  alt="微信公众号"
                />
                <span>公众号</span>
              </div>
              <div class="qrcode-card ml-4">
                <img
                  src="@/assets/bb-helper-wechat-qrcode.webp"
                  alt="BB_小助手"
                />
                <span>BB 小助手</span>
              </div>
            </div>
          </div>
          <div v-else class="w-1/2 pt-2">
            <a href="https://discord.gg/huyw7gRsyA" target="_blank">
              <img
                src="https://discordapp.com/api/guilds/861117579216420874/widget.png?style=banner4"
                alt="Discord Invite"
              />
            </a>
          </div>
        </div>
      </template>
    </DrawerContent>
  </Drawer>
</template>

<script lang="ts" setup>
import type { Node, Tag } from "@markdoc/markdoc";
import { storeToRefs } from "pinia";
import { ref, reactive, watch, computed } from "vue";
import { Drawer, DrawerContent } from "@/components/v2";
import { useLanguage } from "@/composables/useLanguage";
import { useUIStateStore, useHelpStore, usePageMode } from "@/store";

const [
  { default: Markdoc },
  { markdocConfig, DOMPurifyConfig },
  { default: yaml },
  { default: DOMPurify },
] = await Promise.all([
  import("@markdoc/markdoc"),
  import("./config"),
  import("js-yaml"),
  import("dompurify"),
]);

interface State {
  frontmatter: Record<string, string>;
  html: string;
}

const active = ref(false);
const { locale } = useLanguage();
const uiStateStore = useUIStateStore();
const helpStore = useHelpStore();
const helpStoreState = storeToRefs(helpStore);
const pageMode = usePageMode();

const shouldShowHelpDrawer = computed(() => {
  return pageMode.value === "BUNDLED";
});
const helpId = computed(() => helpStoreState.currHelpId.value);
const isGuide = computed(() => helpStoreState.openByDefault.value);

const state = reactive<State>({
  frontmatter: {},
  html: "",
});

watch(helpId, async (id) => {
  if (id) {
    const res = await fetch(
      `/help/${
        locale.value === "zh-CN" ? "zh" : locale.value === "ja-JP" ? "ja" : "en"
      }/${id}.md`
    );
    const markdown = await res.text();
    const ast: Node = Markdoc.parse(markdown);
    const content = Markdoc.transform(ast, markdocConfig) as Tag;

    content.attributes.class = "markdown-body"; // style help content
    const html: string = Markdoc.renderers.html(content);

    state.frontmatter = ast.attributes.frontmatter
      ? (yaml.load(ast.attributes.frontmatter) as Record<string, string>)
      : {};
    state.html = DOMPurify.sanitize(html, DOMPurifyConfig);
    activate();
  } else {
    state.frontmatter = {};
    state.html = "";
    deactivate();
  }
});

const onClose = () => {
  if (isGuide.value) {
    if (!uiStateStore.getIntroStateByKey(`${helpId.value}`)) {
      uiStateStore.saveIntroStateByKey({
        key: `${helpId.value}`,
        newState: true,
      });
    }
  }
  helpStore.exitHelp();
};

const activate = () => (active.value = true);

const deactivate = () => (active.value = false);
</script>

<style scoped>
.qrcode-card {
  @apply w-20 flex flex-col items-center justify-start text-xs;
}
</style>
