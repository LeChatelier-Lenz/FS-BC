
<template>
  <div class="user-info" >
    <div class="info-detail">
      <div class="basic">
          <div class="block">
            <el-avatar :size="50" :src="circleUrl" />
          </div>
          <p>邮箱：{{this.userInfo.useremail}}</p>
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
        useremail:"826733088@qq.com",
        username: "Alice",
        user_id: "001"
      },
      balance: 0.00,
      circleUrl: "https://avatars.githubusercontent.com/u/47231169?v=4",
      contractList: [
        {
          "business_id": "Loan1725414531",
          "amount": 20000.0,
          "issuer": "Alice",
          "created_at": "2024-09-01 14:23:19",
          "updated_at": "2024-09-01 15:11:22",
          "state": "Approved"
        },
        {
          "business_id":"Insurance1725103122",
          "amount": 1000,
          "issuer": "Alice",
          "created_at": "2024-09-01 00:52:02",
          "updated_at": "2024-09-01 02:15:33",
          "state": "Rejected"
        },
        {
          "business_id":"Loan1724118531",
          "amount": 50000,
          "issuer": "Alice",
          "created_at":"2024-08-20 09:48:51",
          "updated_at":"2024-08-29 18:28:29",
          "state":"Claimed"
        }
      ]
    };
  },
  mounted() {
    console.log("UserInfo mounted");
    // this.GetAllContractInfo();
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

.contract-history p {
  margin: 0;
  font-size: 14px;
  color: #8B8E98;
  text-overflow: ellipsis;
  font-family:"Arial Unicode MS";
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
  width: 200px;
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
