
<template>
  <div class="user-info" >
    <div class="info-detail">
      <div class="basic">
          <div class="block">
            <el-avatar :size="50" :src="circleUrl" />
          </div>
          <p>用户名：{{this.userInfo.username}}</p>
          <p>用户ID：{{this.userInfo.user_id}}</p>
      </div>
      <div class="asset">
        <p>总资产（金额）：{{this.balance}} ￥</p>
      </div>
    </div>
    <div class="contract-history">
      <el-scrollbar wrap-class="scrollbar-wrapper" style="height: 350px">
        <!--  //timeline组件，根据contractList数据循环渲染 -->
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in contractList"
            :key="index"
            :timestamp="item.created_at"
            placement="top"
          >
            <el-card>
              <p>合同编号: {{ item.business_id }}</p>
              <p>合同金额: {{ item.amount }}</p>
              <p>发行方: {{ item.issuer }}</p>
              <p>状态: {{ item.state }}</p>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </el-scrollbar>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      userInfo:{
        username: "Alice",
        user_id: "001"
      },
      balance: 0,
      circleUrl: "https://avatars.githubusercontent.com/u/47231169?v=4",
      contractList: [
        {
          "business_id": "string",
          "amount": 0,
          "issuer": "string",
          "created_at": "string",
          "updated_at": "string",
          "state": "string"
        },
        {
          "business_id":"x001",
          "amount": 1000,
          "issuer": "Alice",
          "created_at": "2021-09-01",
          "updated_at": "2021-09-01",
          "state": "active"
        },
        {
          "business_id":"x002",
          "amount": 2000,
          "issuer": "Bob",
          "created_at": "2021-09-02",
          "updated_at": "2021-09-02",
          "state": "active"
        }
      ]
    };
  },
  mounted() {
    console.log("UserInfo mounted");
    this.GetAllContractInfo();
  },
  methods: {
    GetAllContractInfo() {
      // 获取所有合同信息
      axios.get("http://127.0.0.1:8001/ecosys/contract",
          {
            params: {
              user_id: "001"
            }
          }
      ).then(response => {
            this.contractList = [];
            console.log(response.data);
            //从response.data中遍历填充contractList
            for (let i = 0; i < response.data.length; i++) {
              this.contractList.push(response.data[i]);
            }
          }
      ).catch(
          error => {
            console.log(error);
            alert("获取合同信息失败");
          }
      );
    }
  }
};
</script>

<style scoped>
.user-info {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.contract-history {
  padding: 10px;
  border: 4px solid #ececec;
  border-radius: 10px;
  margin-top: 2rem;
}
.info-detail {
  padding: 10px;
  display: flex;
  flex-direction: row;
  justify-content:flex-start;
  margin-top: 2rem;
  margin-left :20px;
}

.basic {
  width: 150px;
  margin-left: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-right: 40px;
}

.basic p {
  margin: 0;
  font-size: 14px;
  color: #8B8E98;
  text-overflow: ellipsis;
}

.asset {
  width: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.asset p {
  width: 200px;
  margin: 20px;
  font-size: 17px;
  text-overflow: ellipsis;
}

</style>
