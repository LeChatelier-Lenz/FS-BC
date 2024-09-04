<script setup lang="ts">

import {ElIcon, ElMenu, ElMenuItem, ElRadioButton, ElRadioGroup, ElSubMenu} from "element-plus";
import {CreditCard, DocumentAdd, House, Stamp, Wallet} from "@element-plus/icons-vue";
import {ref} from "vue";
const isCollapse = ref(true)
const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}

</script>

<template>
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
</template>

<style scoped>

.RouterLink {
  margin: 0;
  padding: 0;
}

.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 230px;
  min-height: 400px;
}

</style>
