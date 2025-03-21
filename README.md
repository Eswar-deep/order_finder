# 🚀 Order Finder 🗺️📍

A **serverless, geolocation-based order tracking system** built using **AWS Lambda, DynamoDB, API Gateway, and Google Maps**.  
Customers can place orders, and admins can view all orders on an **interactive map dashboard**.

## 🛠️ Tech Stack
- **Backend:** Go (AWS Lambda)
- **Database:** DynamoDB
- **API:** API Gateway (Serverless Framework)
- **Frontend:** HTML + JavaScript (S3 Hosted)
- **Geolocation:** Google Maps API

---

## 🌍 **Live Demo**
🔹 **Order Form (Customers Place Orders):**  
👉 [Order Now](http://order-finder-frontend.s3-website-us-east-1.amazonaws.com/index.html)  

🔹 **Admin Order Map (View All Orders):**  
👉 [View Orders](http://order-finder-frontend.s3-website-us-east-1.amazonaws.com/map.html)  

---

## 🖼️ Screenshots

### 📝 Customer Order Form
![Order Form](https://raw.githubusercontent.com/Eswar-deep/order_finder/main/Screenshot%202025-03-21%20053402.png)

### 🌍 Order Map Dashboard
![Order Map](https://raw.githubusercontent.com/Eswar-deep/order_finder/main/Screenshot%202025-03-21%20053429.png)

---

## 🏗️ **Deployment Instructions**

### 1️⃣ **Setup AWS & Serverless**
Make sure you have **AWS CLI** and **Serverless Framework** installed:

```bash
npm install -g serverless
aws configure
