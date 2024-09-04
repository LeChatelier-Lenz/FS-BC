<script setup lang="ts">
import UserInfo from "@/components/UserInfo.vue";
</script>

<template>
  <main>
    <div class="content">
    <UserInfo />
    </div>
  </main>
</template>

<style scoped>

.content{
  margin-left: 80px;
  height: 100%;
  width: 1000px;
  background: #ffffff;
  box-shadow:  20px 20px 60px #d9d9d9,
  -20px -20px 60px #ffffff;
}

</style>
