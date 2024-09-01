<template>
  <div class="modal">
    <div class="form">
      <h2>贷款合同履约检查</h2>
      <div class="separator">
        <hr class="line">
        <p> contract-info </p>
        <hr class="line">
      </div>
      <form @submit.prevent="handleSubmit" class="sheet">
        <div class="input_container">
          <label for="user_id" class="input_label">用户ID</label>
          <input id="user_id" v-model="form.user_id" class="input_field" required>
        </div>
        <div class="input_container">
          <div class="split">
            <div>
              <label for="business_id" class="input_label">业务ID</label>
              <input id="business_id" v-model="form.business_id" class="input_field" required>
            </div>
            <button type="button" class="query-btn" @click="fetchContractInfo">查询</button>
          </div>
        </div>
        <div class="input_container">
          <label class="input_label">信用评分</label>
          <input id="credit" type="number" v-model="form.conditions.credit" class="input_field" required>
        </div>
        <div class="input_container">
          <label class="input_label">收入</label>
          <input id="income" type="number" v-model="form.conditions.income" class="input_field" required>
        </div>
<!--        <div class="input_container">-->
<!--          <label for="current_time" class="input_label">当前时间</label>-->
<!--          <input id="current_time" type="datetime-local" v-model="form.current_time" class="input_field" required>-->
<!--        </div>-->
        <button type="submit" class="purchase--btn">检查履约</button>
      </form>
    </div>
    <div class="contract-card">
      <h3>合同信息</h3>
      <p v-if="contractInfo">
        <strong>合同状态:</strong> {{ contractInfo.status }}<br>
        <strong>合同金额:</strong> {{ contractInfo.amount }}<br>
        <strong>合同发行方:</strong> {{ contractInfo.issuer }}<br>
        <strong>利率:</strong> {{ contractInfo.rate }}%
      </p>
      <p v-else style="color: #8B8E98">
        点击查询按钮可获取合同信息
      </p>
    </div>
  </div>
</template>

<script>
//从assets文件夹导入css文件
import "../assets/sheet.css";
import axios from "axios";
export default {
  data() {
    return {
      form: {
        user_id: '',
        business_id: '',
        conditions: {
          credit: 0,
          income: 0
        },
        current_time: ''
      },
      contractInfo: null // 合同信息
    };
  },
  methods: {
    handleSubmit() {
      // 处理提交逻辑
      console.log(this.form);
      // 发送合同启动请求
      // 利用axios发送请求
      // 替换为实际的API调用
      this.form.current_time = new Date().getSeconds();
      axios.post('http://127.0.0.1:8001/ecosys/loan/check_contract', this.form).then(
          response => {
            console.log(response.data);
            //提示窗
            alert("检查成功,贷款已强制还款");
          }
      ).catch(
          error => {
            console.log(error);
            //提示窗
            alert("检查失败"+error);
          }
      )
    },
    async fetchContractInfo() {
      // 模拟从服务器获取合同信息
      // 替换为实际的API调用
      axios.get('http://localhost:8000/ecosys/query_contract', {
        params: {
          user_id: this.form.user_id,
          business_id: this.form.business_id,
          type: "loan"
        }
      }).then(response => {
        console.log(response.data);
        this.contractInfo = {
          status: response.data.state,
          amount: response.data.amount,
          issuer: response.data.issuer,
          rate: response.data.rate
        }
      }).catch(error => {
        console.error('获取合同信息失败', error);
        alert("获取合同信息失败");
      });
    }
  }
};
</script>

