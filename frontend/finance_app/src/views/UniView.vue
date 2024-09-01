<template>
  <div class="modal">
    <div class="form">
      <h2>合同申请</h2>
      <el-form :model="form" @submit.native.prevent="handleSubmit">
        <el-form-item  class="input_container">
          <label class="input_label">用户ID</label>
          <el-input v-model="form.user_id" class="input_field" required></el-input>
        </el-form-item>
        <el-form-item  class="input_container">
          <label class="input_label">密码</label>
          <el-input type="password" v-model="form.password" class="input_field" required></el-input>
        </el-form-item>
        <el-form-item class="input_container">
          <label class="input_label">业务类型</label>
          <div class="radio_group">
            <el-radio-group v-model="form.business_type" fill="black">
              <el-radio-button label="Insurance">保险</el-radio-button>
              <el-radio-button label="Loan">贷款</el-radio-button>
            </el-radio-group>
          </div>
        </el-form-item>
        <el-form-item  class="input_container">
          <label class="input_label">金额</label>
          <el-input
              v-model="form.amount"
              class="input_field"
              :formatter="(value) => `￥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
              :parser="(value) => value.replace(/\$\s?|(,*)/g, '')"
              required></el-input>
        </el-form-item>
        <el-form-item  class="input_container">
          <label class="input_label">利率</label>
<!--          // 设置type为number，step为0.001，formatter为`${value} %`，使输入框只能输入数字，且保留三位小数-->
          <el-input
              type="number"
              step="0.01"
              v-model="form.rate"
              class="input_field"
              required><template #append>%</template>
              </el-input>
        </el-form-item>
        <el-form-item class="input_container">
          <label class="input_label">发行商</label>
          <div class="radio_group">
            <el-radio-group v-model="form.issuer">
              <el-radio label="issuer1">发行商1</el-radio>
              <el-radio label="issuer2">发行商2</el-radio>
              <el-radio label="issuer3">发行商3</el-radio>
            </el-radio-group>
          </div>
        </el-form-item>
        <el-form-item class="input_container">
          <label class="input_label">期限</label>
          <el-input type="number" v-model="form.period"  class="input_field" required> <template #append>天</template></el-input>
        </el-form-item>
        <el-form-item>
          <button type="submit" class="purchase--btn">创建合同</button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import axios  from "axios";
export default {
  data() {
    return {
      form: {
        user_id: '',
        business_id: '',
        password: '',
        business_type: '',
        amount: 0,
        rate: 0,
        issuer: '',
        period: 0
      }
    };
  },
  methods: {
    handleSubmit() {
      // 处理提交逻辑
      console.log(this.form);
      this.business_id = this.form.business_type + new Date().getSeconds();
      this.form.amount = parseFloat(this.form.amount);
      this.form.rate = this.form.rate / 100;
      this.form.period = parseInt(this.form.period);
      this.form.current_time = new Date().getSeconds();
      axios.post('http://127.0.0.1:8001/ecosys/new_contract', this.form).then(
        response => {
          console.log(response.data);
          //提示窗
          alert("合同创建成功");
        }
      ).catch(
        error => {
          console.log(error);
          //提示窗
          alert("合同创建失败");
        }
      );
    }
  }
};
</script>

<style scoped>
.modal {
  margin: 0 auto;
  max-width: 90vw; /* 使宽度适应视口 */
  width: 500px; /* 最大宽度 */
  background: #FFFFFF;
  box-shadow: 0px 187px 75px rgba(0, 0, 0, 0.01), 0px 105px 63px rgba(0, 0, 0, 0.05), 0px 47px 47px rgba(0, 0, 0, 0.09), 0px 12px 26px rgba(0, 0, 0, 0.1), 0px 0px 0px rgba(0, 0, 0, 0.1);
  border-radius: 26px;
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
}

.input_label {
  font-size: 13px;
  color: #8B8E98;
  font-weight: 600;
  padding-right: 10px;
}

.input_field {
  width: 80%; /* 使输入字段宽度适应容器 */
  height: 30px;
  padding: 0 0 0 0px;
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

.radio_group {
  display: flex;
  gap: 10px;
}

.radio_group label {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #8B8E98;
}

.purchase--btn {
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

.purchase--btn:hover {
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
