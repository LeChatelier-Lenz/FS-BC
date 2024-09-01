
<template>
  <main>
    <div class="modal">

      <form class="form" @submit.prevent="handleSubmit">
        <h2>支付操作</h2>
        <div class="separator">
          <hr class="line">
          <p> contract-info </p>
          <hr class="line">
        </div>
        <div class="input_container">
          <label for="user_id_field" class="input_label">用户ID</label>
          <input id="user_id_field" class="input_field" v-model="this.user_id"  type="text" name="user_id" title="用户ID" placeholder="请输入用户ID" required>
        </div>
        <div class="input_container">
          <label for="password_field" class="input_label">密码</label>
          <input id="password_field" class="input_field" v-model="this.password" type="password" name="password" title="密码" placeholder="请输入密码" required>
        </div>
        <div class="input_container">
          <label for="amount_field" class="input_label">金额</label>
          <input id="amount_field" class="input_field" v-model="this.amount" type="number" name="amount" title="金额" placeholder="请输入金额" required>
        </div>
        <div class="input_container">
          <label for="type_field" class="input_label">类型</label>
          <select id="type_field" class="input_field" name="type" v-model="this.type" title="类型">
            <option value="deposit">存款</option>
            <option value="transfer">转账</option>
          </select>
        </div>
        <div class="input_container" v-if="this.type === 'transfer'">
          <label for="target_user_id_field" class="input_label">目标用户ID</label>
          <input id="target_user_id_field" class="input_field" v-model="this.target_user_id" type="text" name="target_user_id" title="目标用户ID" placeholder="请输入目标用户ID">
        </div>
        <button class="submit--btn">提交</button>
      </form>
    </div>
  </main>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      type: 'deposit',
      user_id: '',
      password: '',
      amount: 0,
      target_user_id: '',
      submitTime: 0
    }
  },
  methods: {
    handleSubmit() {
      console.log('submit')
      this.submitTime = new Date().getSeconds();
      if (this.type === 'transfer') {
        console.log('transfer')
        axios.post("http://127.0.0.1:8001/ecosys/pay/transfer",
            {
              user_id: this.user_id,
              password: this.password,
              amount: this.amount,
              target_user_id: this.target_user_id,
              current_time: this.submitTime
            }).then(response => {
          console.log(response.data)
          //提交成功弹窗
          alert("转账成功");
        }).catch(error => {
          console.log(error)
          //提交失败弹窗
          alert("转账失败");
        })
      } else {
        console.log('deposit')
        axios.post("http://127.0.0.1:8001/ecosys/pay/deposit",
            {
              user_id: this.user_id,
              password: this.password,
              amount: this.amount,
              current_time: this.submitTime
            }).then(response => {
          console.log(response.data)
          //提交成功弹窗
          alert("存款成功");
        }).catch(error => {
          console.log(error)
          //提交失败弹窗
          alert("存款失败");
        })
      }
    }
  }
}
</script>


<style scoped>
.modal {
  margin: 0 auto;
  width: fit-content;
  height: fit-content;
  background: #FFFFFF;
  box-shadow: 0px 187px 75px rgba(0, 0, 0, 0.01), 0px 105px 63px rgba(0, 0, 0, 0.05), 0px 47px 47px rgba(0, 0, 0, 0.09), 0px 12px 26px rgba(0, 0, 0, 0.1), 0px 0px 0px rgba(0, 0, 0, 0.1);
  border-radius: 26px;
  max-width: 450px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px;
}

.input_container {
  width: 100%;
  height: fit-content;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.input_label {
  font-size: 10px;
  color: #8B8E98;
  font-weight: 600;
}

.input_field {
  width: auto;
  height: 40px;
  padding: 0 0 0 16px;
  border-radius: 9px;
  outline: none;
  background-color: #F2F2F2;
  border: 1px solid #e5e5e500;
  transition: all 0.3s cubic-bezier(0.15, 0.83, 0.66, 1);
}

.input_field:focus {
  border: 1px solid transparent;
  box-shadow: 0px 0px 0px 2px #242424;
  background-color: transparent;
}

.submit--btn {
  height: 55px;
  background: #F2F2F2;
  border-radius: 11px;
  border: 0;
  outline: none;
  color: #ffffff;
  font-size: 13px;
  font-weight: 700;
  background: linear-gradient(180deg, #363636 0%, #1B1B1B 50%, #000000 100%);
  box-shadow: 0px 0px 0px 0px #FFFFFF, 0px 0px 0px 0px #000000;
  transition: all 0.3s cubic-bezier(0.15, 0.83, 0.66, 1);
}

.submit--btn:hover {
  box-shadow: 0px 0px 0px 2px #FFFFFF, 0px 0px 0px 4px #0000003a;
}

/* Reset input number styles */
.input_field::-webkit-outer-spin-button,
.input_field::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.input_field[type=number] {
  -moz-appearance: textfield;
}
</style>
```

这段代码实现了一个包含用户ID、业务ID、密码、金额和类型（转账/存款）的表单。如果选择类型为转账，则会额外显示一个目标用户ID的输入框。页面文字已改为中文。
