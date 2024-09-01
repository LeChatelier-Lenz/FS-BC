<script setup lang="ts">
import { ref } from 'vue'
import { ElMenu, ElMenuItem, ElMenuItemGroup, ElSubMenu, ElRadioButton, ElRadioGroup,ElIcon,ElDivider } from 'element-plus'
import {
  House,
  DocumentAdd,
  Stamp,
  Wallet,
  CreditCard
} from '@element-plus/icons-vue'


const isCollapse = ref(true)
const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
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
      <h1>Finance App</h1>
    </div>
    <div class="header-right">
      <el-button type="primary" icon="el-icon-share" circle></el-button>
      <el-button type="danger" icon="el-icon-switch-button" circle></el-button>
    </div>
  </div>
  <el-divider></el-divider>
  <div class = "body">
  <div class = "sidebar">
  <el-radio-group v-model="isCollapse" style="margin-bottom: 20px" fill="black">
    <el-radio-button :value="true"> < </el-radio-button>
    <el-radio-button :value="false"> > </el-radio-button>
  </el-radio-group>
  <el-menu
      default-active="2"
      class="el-menu-vertical-demo"
      :collapse="isCollapse"
      @open="handleOpen"
      @close="handleClose"
  >

    <RouterLink class="RouterLink" to="/home">
    <el-menu-item index="1">
      <el-icon><House /></el-icon>
      <template #title><span>系统主页/SYSTEM</span></template>
    </el-menu-item>
    </RouterLink>

    <RouterLink class="RouterLink" to="/ecosys">
      <el-menu-item index="2">
        <el-icon><DocumentAdd /></el-icon>
        <template #title><span>新建业务/CONTRACT</span></template>
      </el-menu-item>
    </RouterLink>

    <el-sub-menu index="3">
      <template #title>
        <el-icon><Stamp /></el-icon>
        <span>保险窗口/INSURANCE</span>
      </template>
      <RouterLink class="RouterLink" to="/ecosys/insuranceStart">
        <el-menu-item index="3-1">保险金存入</el-menu-item>
      </RouterLink>
      <RouterLink class="RouterLink" to="/ecosys/insuranceCheck">
        <el-menu-item index="3-2">合同履约检查</el-menu-item>
      </RouterLink>
    </el-sub-menu>

    <el-sub-menu index="4">
      <template #title>
        <el-icon><Wallet /></el-icon>
        <span>贷款窗口/LOAN</span>
      </template>
      <RouterLink class="RouterLink" to="/ecosys/loanStart">
        <el-menu-item index="4-1"> 贷款发放 </el-menu-item>
      </RouterLink>
      <RouterLink class="RouterLink" to="/ecosys/loanCheck">
        <el-menu-item index="4-2">合同履约检查</el-menu-item>
      </RouterLink>
    </el-sub-menu>

    <el-sub-menu index="5">
      <template #title>
        <el-icon><CreditCard /></el-icon>
        <span>支付窗口/PAYMENT</span>
      </template>
      <RouterLink class="RouterLink" to="/ecosys/pay">
        <el-menu-item index="5-1"> 转账/存款 </el-menu-item>
      </RouterLink>
    </el-sub-menu>


  </el-menu>
  </div>
  <div class = "content">
    <RouterView />
  </div>
  </div>
<!--  <RouterView />-->
</template>

<style scoped>
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 230px;
  min-height: 400px;
}
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
}

.el-menu-vertical-demo span {
  font-weight: bold;
}



.RouterLink {
  margin: 0;
  padding: 0;
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
