<script setup lang="ts">
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import InfiniteLoading from "v3-infinite-loading";
import "v3-infinite-loading/lib/style.css";
import { getCardsPage, Card } from "@/api/card";
import { uploadAvatar } from "@/api/site";
import CalHeatmap from "@/components/CalHeatmap.vue";

const router = useRouter();
const route = useRoute();

const cards = ref<Card[]>([]);
let page = 1;
let pageSize = 20;
const load = async ($state) => {
  try {
    const resp = await getCardsPage(
      page,
      route.query.q as string,
      route.query.d as string
    );
    cards.value.push(...resp.data.cards);
    if (resp.data.cards.length < pageSize) $state.complete();
    else {
      $state.loaded();
    }
    page++;
  } catch (error) {
    $state.error();
  }
};

const isLogin = ref(false);
if (localStorage.getItem("accessToken")) {
  isLogin.value = true;
}

const jumpCardDetailPage = (id: number) => {
  let routerName = "card-detail";
  if (isLogin.value) {
    routerName = "user-card-detail";
  }
  router.push({ name: routerName, params: { id: id } });
};

const jumpPostPage = () => {
  router.push({ name: "user-card-post", params: { id: 0 } });
};

const jumpDayCardPage = async (date: string) => {
  router.push({ name: "card-page", query: { d: date } });
};

const jumpCardStatPage = async (date: string) => {
  router.push({ name: "card-stat" });
};

let siteInfo = JSON.parse(localStorage.getItem('siteInfo') || '{}');
const nickname = siteInfo.nickname;
const avatar = siteInfo.avatar;

const tryToUploadAvatar = async (e: any) => {
  const file = e.target.files[0];
  const formData = new FormData();
  formData.append("file", file);
  const res = await uploadAvatar(formData);
  if (res.code === 200) {
    // 上传成功刷新页面
    location.reload();
  }
};

const handleAvatarClick = () => {
  if (!isLogin.value) {
    return;
  }
  const input = document.getElementById("upload") as HTMLInputElement;
  input.click();
};
</script>

<template>
  <div class="card-list-bg">
    <div class="card-list">
      <div class="card-list-item">
        <div class="ribbon-header" @click="jumpCardStatPage">
          <div class="ribbon">
            <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M5.81836 6.72729V14H13.0911" stroke="#fffefe" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/><path d="M4 24C4 35.0457 12.9543 44 24 44V44C35.0457 44 44 35.0457 44 24C44 12.9543 35.0457 4 24 4C16.598 4 10.1351 8.02111 6.67677 13.9981" stroke="#fffefe" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/><path d="M24.005 12L24.0038 24.0088L32.4832 32.4882" stroke="#fffefe" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </div>
        </div>
        <CalHeatmap @clickBox="jumpDayCardPage" />
        <hr />
        <div class="card-list-item-user-info">
          <img :src="avatar" @click="handleAvatarClick"/>
          <input id="upload" type="file" accept="image/png, image/jpeg" @change="tryToUploadAvatar" style="display: none" />
          <div class="card-list-item-user-nickname">
            <p>{{ nickname }}</p>
          </div>
        </div>
      </div>
      <div class="card-list-item" v-for="card in cards" :key="card.id">
        <div class="card-content" v-html="card.content" v-highlight></div>
        <hr />
        <div class="card-footer">
          <span>📅 {{ card.created_at }} <span v-if="isLogin">👀 {{ card.pv }}</span></span>
          <span @click="jumpCardDetailPage(card.id)">详情</span>
        </div>
      </div>
      <InfiniteLoading
        @infinite="load"
        :slots="{ complete: '没有更多卡片了', error: '加载失败' }"
      />
    </div>
    <button v-if="isLogin" class="post-card-btn" @click="jumpPostPage()">+</button>
  </div>
</template>

<style scoped>
.card-list-bg {
  display: flex;
  flex-direction: row;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(-45deg, #ee7752, #e73c7e, #23a6d5, #23d5ab);
  background-size: 400% 400%;
}

.card-list {
  padding: 20px 20px;
  margin-top: 60px;
  max-width: 1200px;
  width: calc(100% - 40px);
}
.card-list-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
  font-weight: 300;
  /*width: 100%;*/
  padding: 10px;
  margin: 0 auto 10px auto;
  text-align: center;
  background: #ffffff;
  border-radius: 10px;
}

.card-list-item hr {
  display: flex;
  position: relative;
  margin: 8px 0 0 0;
  border: 1px dashed #4259ef23;
  width: 100%;
}

.card-footer {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: baseline;
}

:deep(.card-content) {
  height: 100%;
  width: 100%;
  word-wrap: break-word;
  text-align: left;
}

.post-card-btn {
  position: fixed;
  bottom: 80px;
  right: 30px;
  border-radius: 50%;
  border: none;
  width: 60px;
  height: 60px;
  font-size: 40px;
  background-color: #1fc6e6;
  color: #f7f9fe;
  box-shadow: 10px 10px 20px rgba(0, 0, 0, 0.3),
    -10px -10px 20px rgba(255, 255, 255, 0.3);
}

.ribbon-header {
  cursor: pointer;
  float: right;
  display: flex;
  justify-content: right;
  align-items: center;
  height: 0;
  margin-top: 10px;
  margin-right: -10px;
}

.ribbon {
  font-size: 14px;
  color: #fff;
  display: flex;
  justify-content: right;
  align-items: center;
  padding: 5px 0 5px 0;
}
.ribbon {
  --r: 0.8em; /* control the cutout */

  border-block: 0.5em solid #0000;
  padding-inline: calc(var(--r) + 0.25em) 0.5em;
  line-height: 1.8;
  clip-path: polygon(
    0 0,
    100% 0,
    100% 100%,
    0 100%,
    0 calc(100% - 0.25em),
    var(--r) 50%,
    0 0.25em
  );
  background: rgb(100, 206, 170) padding-box; /* the color  */
  width: fit-content;
}

.card-list-item-user-info {
  display: flex;
  flex-direction: row;
  justify-content: end;
  align-items: center;
}

.card-list-item-user-info>img {
  border-radius: 50%;
  width: 2em;
  height: 2em;
  margin-right: 0.75em;
}

.card-list-item-user-nickname {
  font-size: 1em;
  
}
.card-list-item-user-nickname > p {
  margin: 0;
}
</style>

<style>
.card-content > ul {
  margin-block-start: 0em;
  margin-block-end: 0em;
  padding-inline-start: 20px;
}
</style>
