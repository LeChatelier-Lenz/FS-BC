<script setup lang="ts">

import Navigator from "@/components/Navigator.vue";
import {ref} from "vue";
const isLogged = ref(false);



</script>

<template>
<!--&lt;!&ndash;    <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="125" height="125" />&ndash;&gt;-->

<!--    <div class="wrapper">-->
<!--&lt;!&ndash;      <HelloWorld msg="You did it!" />&ndash;&gt;-->

<!--      <nav>-->
<!--        <RouterLink to="/">Home</RouterLink>-->
<!--        <RouterLink to="/about">About</RouterLink>-->
<!--      </nav>-->
<!--    </div>-->
  <div class = "header">
    <div class="header-left">
      <h1>Trust-Ex</h1>
    </div>
    <div class="header-right">
      <el-button type="primary" icon="el-icon-share" circle></el-button>
      <el-button type="danger" icon="el-icon-switch-button" circle></el-button>
    </div>
  </div>
  <el-divider></el-divider>
  <div class = "body">
    <!-- 当路由为'/'时，不显示sidebar块；当路由不为'/'时，显示sidebar块 -->
    <div class = "sidebar"  v-if="$route.path !== '/'">
      <Navigator />
    </div>
  <div class = "content">
    <RouterView />
  </div>
  </div>
<!--  <RouterView />-->
</template>

<style scoped>

header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}

.body {
  height: 600px;
  width: 100%;
  margin: 0 auto;
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
}

.sidebar {
  position: fixed;
  height: 100%;
  width: 30%;
  margin-left: 0;
  z-index: 100;
}

.sidebar el-menu {
  background-color: white;
}

.content{
  height: 100%;
  position: relative;
  margin: 0 100px 0  150px;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
}

.el-menu-vertical-demo span {
  font-weight: bold;
}







.header {
  position: fixed;
  width: 100%;
  top: 0;
  height: 80px;
  background: #ffffff;
  box-shadow: 10px 10px 30px #d9d9d9, -10px -10px 30px #ffffff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  z-index: 1000;
}

.header-left h1 {
  font-size: 32px; /* 更大的字体 */
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2px;
  margin: 0;
}

.header-right {
  display: flex;
  gap: 15px;
}

</style>
