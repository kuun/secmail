<template>
  <div class="animated-bg">
    <div v-for="i in 4" :key="i" class="bubble"></div>
  </div>
</template>

<style scoped>
.animated-bg {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
  overflow: hidden;
  background: linear-gradient(to bottom, #EEF2FF, #ffffff);
}

.bubble {
  position: absolute;
  background: radial-gradient(circle at center, rgba(59, 130, 246, 0.1) 0%, rgba(59, 130, 246, 0.05) 50%, transparent 70%);
  border-radius: 50%;
}

.bubble:nth-child(1) {
  width: 400px;
  height: 400px;
  left: -100px;
  top: -100px;
  animation: float1 3s infinite linear;
}

.bubble:nth-child(2) {
  width: 300px;
  height: 300px;
  right: -50px;
  top: 30%;
  animation: float2 6s infinite linear;
}

.bubble:nth-child(3) {
  width: 200px;
  height: 200px;
  left: 10%;
  bottom: 10%;
  animation: float3 8s infinite linear;
}

.bubble:nth-child(4) {
  width: 250px;
  height: 250px;
  right: 20%;
  top: 20%;
  animation: float4 3s infinite linear;
}

@keyframes float1 {
  0% { transform: translate(0, 0) rotate(0deg); }
  50% { transform: translate(100px, 100px) rotate(180deg); }
  100% { transform: translate(0, 0) rotate(360deg); }
}

@keyframes float2 {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-50px, 50px) scale(1.1); }
  100% { transform: translate(0, 0) scale(1); }
}

@keyframes float3 {
  0% { transform: translate(0, 0) scale(1.1); }
  50% { transform: translate(100px, -100px) scale(1); }
  100% { transform: translate(0, 0) scale(1.1); }
}

@keyframes float4 {
  0% { transform: translate(0, 0) rotate(0deg) scale(1); }
  50% { transform: translate(-30px, 30px) rotate(-180deg) scale(1.2); }
  100% { transform: translate(0, 0) rotate(-360deg) scale(1); }
}
</style>
