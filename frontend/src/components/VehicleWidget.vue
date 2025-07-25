<template>
  <div class="widget-container" style="--wails-draggable: drag">
    <!-- 苹果风格光点 -->
    <div class="glow-particles">
      <div class="particle particle-1"></div>
      <div class="particle particle-2"></div>
      <div class="particle particle-3"></div>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchData(false)" class="retry-btn">重试</button>
    </div>

    <div v-else-if="vehicleData" class="widget-content">
      <!-- 顶部关闭按钮和时间 -->
      <div class="header drag-region">
        <button class="close-btn no-drag" @click="minimizeWindow">×</button>
        <span class="time">{{ currentTime }}</span>
        <div class="drag-handle" title="拖动窗口">⋮⋮</div>
      </div>

      <!-- 主要内容区域 -->
      <div class="main-content">
        <!-- 左侧信息 -->
        <div class="info-section">
          <div class="stats">
            <span class="range">{{ vehicleData.hmiRidableMile }}km</span>
            <span class="battery">{{ vehicleData.bmssoc }}%</span>
          </div>
          <div class="battery-bar">
            <div
              class="battery-fill"
              :style="{ width: parseInt(vehicleData.bmssoc || 0) + '%' }"
            ></div>
          </div>
        </div>

        <!-- 右侧车辆图片 -->
        <div class="vehicle-section">
          <img
            v-if="vehicleData.vehiclePicUrl"
            :src="vehicleData.vehiclePicUrl"
            alt="Vehicle"
            class="vehicle-image"
            @error="handleImageError"
          />
          <div v-else class="vehicle-placeholder">🛵</div>
        </div>
      </div>

      <!-- 底部控制按钮 -->
      <div class="controls">
        <!-- <button class="control-btn lock-btn">🔒</button>
        <button class="control-btn sound-btn">🔊</button> -->

        <!-- 设置按钮 -->
        <button
          class="control-btn settings-btn"
          @click="openConfigModal"
          title="设置"
        >
          ⚙️
        </button>
        <div class="location-info">
          <div class="location-time">
            📍 {{ vehicleData.location?.locationTime || "位置信息" }}
          </div>
          <div v-if="vehicleData.location?.address" class="location-address">
            {{ vehicleData.location.address }}
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 配置模态框 -->
  <ConfigModal
    :show="showConfigModal"
    @close="closeConfigModal"
    @saved="onConfigSaved"
  />
</template>

<script setup>
import { onMounted, ref, computed } from "vue";
import {
  GetMockVehicleData,
  GetVehicleData,
  GetConfig,
  MoveToCorner,
  MinimizeToTray,
} from "../../wailsjs/go/main/App";
import ConfigModal from "./ConfigModal.vue";

const vehicleData = ref(null);
const loading = ref(false);
const error = ref(null);
const showConfigModal = ref(false); 
const hasConfig = ref(false);

const currentTime = computed(() => {
  const now = new Date();
  return now.toLocaleTimeString("zh-CN", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
});

const fetchData = async (useMock = false) => {
  loading.value = true;
  error.value = null;

  try {
    let data;
    try {
      data = await GetVehicleData();
      console.log("GetVehicleData", data);
      vehicleData.value = data;
      error.value = null;
    } catch (apiErr) {
      data = await GetMockVehicleData();
      vehicleData.value = data;
      error.value = null; // 使用模拟数据时不显示错误
    }
  } catch (err) {
    error.value = "获取数据失败: " + err.message;
    console.error("Failed to fetch vehicle data:", err);
  } finally {
    loading.value = false;
  }
};

const checkConfig = async () => {
  try {
    const config = await GetConfig();
    hasConfig.value = config && config.token && config.vehicleId;
    return hasConfig.value;
  } catch (err) {
    console.error("检查配置失败:", err);
    hasConfig.value = false;
    return false;
  }
};

const openConfigModal = () => {
  showConfigModal.value = true;
};

const closeConfigModal = () => {
  showConfigModal.value = false;
};

const onConfigSaved = async () => {
  showConfigModal.value = false;
  await checkConfig();
  if (hasConfig.value) {
    await fetchData(false);
  }
};

const handleImageError = (event) => {
  event.target.style.display = "none";
};

const minimizeWindow = async () => {
  try {
    await MinimizeToTray();
  } catch (err) {
    console.error("最小化失败:", err);
  }
};

onMounted(async () => {
  const configExists = await checkConfig();
  if (configExists) {
    await fetchData(false);
  } else {
    // 直接显示模拟数据，不强制要求配置
    await fetchData(true);
  }
});
</script>

<style scoped>
.widget-container {
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 20px;
  padding: 12px;
  box-sizing: border-box;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  position: relative;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  overflow: hidden;
}

/* 纯净白色系光晕效果 */
.widget-container::before {
  content: "";
  position: absolute;
  top: -100px;
  left: -100px;
  width: calc(100% + 200px);
  height: calc(100% + 200px);
  background: radial-gradient(
      circle at 30% 20%,
      rgba(255, 255, 255, 0.4) 0%,
      rgba(255, 255, 255, 0.2) 30%,
      transparent 60%
    ),
    radial-gradient(
      circle at 70% 30%,
      rgba(248, 250, 252, 0.35) 0%,
      rgba(241, 245, 249, 0.15) 40%,
      transparent 70%
    ),
    radial-gradient(
      circle at 50% 80%,
      rgba(255, 255, 255, 0.3) 0%,
      rgba(248, 250, 252, 0.1) 35%,
      transparent 65%
    ),
    radial-gradient(
      circle at 85% 60%,
      rgba(241, 245, 249, 0.25) 0%,
      rgba(255, 255, 255, 0.1) 30%,
      transparent 50%
    );
  animation: whiteAurora 18s ease-in-out infinite;
  pointer-events: none;
  z-index: -1;
}

/* 额外的光晕层 */
.widget-container::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.1) 0%,
    transparent 50%,
    rgba(255, 255, 255, 0.05) 100%
  );
  border-radius: 20px;
  pointer-events: none;
  z-index: 1;
  mix-blend-mode: overlay;
}

@keyframes whiteAurora {
  0%,
  100% {
    transform: rotate(0deg) scale(1);
    opacity: 0.8;
  }
  33% {
    transform: rotate(120deg) scale(1.02);
    opacity: 0.6;
  }
  66% {
    transform: rotate(240deg) scale(0.98);
    opacity: 0.9;
  }
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #ddd;
  border-top: 2px solid #007aff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #ff3b30;
  font-size: 11px;
  text-align: center;
}

.retry-btn {
  background: #007aff;
  color: white;
  border: none;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 10px;
  margin-top: 8px;
  cursor: pointer;
}

.widget-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 3;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  cursor: move;
  user-select: none;
  padding: 4px 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 16px;
  color: #666;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  z-index: 10;
  position: relative;
}

.close-btn:hover {
  color: #ff3b30;
  background: rgba(255, 59, 48, 0.1);
  border-radius: 50%;
}

.time {
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.drag-handle {
  font-size: 14px;
  color: #999;
  cursor: move;
  padding: 2px 4px;
  border-radius: 4px;
  transition: color 0.2s;
}

.drag-handle:hover {
  color: #666;
  background: rgba(255, 255, 255, 0.1);
}

.main-content {
  flex: 1;
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.info-section {
  flex: 1;
}

.stats {
  display: flex;
  gap: 12px;
  margin-bottom: 6px;
}

.range,
.battery {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.battery-bar {
  width: 100%;
  height: 4px;
  background: #e0e0e0;
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 12px;
}

.battery-fill {
  height: 100%;
  background: linear-gradient(90deg, #00d4ff 0%, #0099cc 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.vehicle-section {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.vehicle-image {
  max-width: 280px;
  max-height: 160px;
  object-fit: contain;
}

.vehicle-placeholder {
  font-size: 40px;
  opacity: 0.6;
}

.controls {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.control-btn {
  width: 32px;
  height: 32px;
  background: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.control-btn:hover {
  background: #f0f0f0;
}

.location-info {
  flex: 1;
  font-size: 10px;
  color: #007aff;
  text-align: right;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: center;
  gap: 2px;
}

.location-time {
  font-size: 12px;
  color: #666;
}

.location-address {
  font-size: 12px;
  color: #666;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.settings-btn:hover {
  opacity: 1;
}

.position-menu {
  position: absolute;
  bottom: 50px;
  left: 12px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 8px;
  z-index: 1000;
  min-width: 80px;
}

.position-option {
  padding: 8px 12px;
  font-size: 12px;
  color: #333;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.position-option:hover {
  background: #f0f0f0;
}

.minimize-btn {
  color: white;
}

.position-btn {
  color: white;
}

.glow-particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 2;
}

.particle {
  position: absolute;
  border-radius: 50%;
  background: radial-gradient(
    circle,
    rgba(255, 255, 255, 0.9) 0%,
    rgba(248, 250, 252, 0.5) 40%,
    transparent 70%
  );
  filter: blur(1px);
  box-shadow: 0 0 6px rgba(255, 255, 255, 0.4);
}

.particle-1 {
  width: 4px;
  height: 4px;
  top: 25%;
  left: 15%;
  animation: sparkle1 8s ease-in-out infinite;
}

.particle-2 {
  width: 3px;
  height: 3px;
  top: 70%;
  right: 20%;
  animation: sparkle2 6s ease-in-out infinite 2s;
}

.particle-3 {
  width: 5px;
  height: 5px;
  top: 45%;
  right: 35%;
  animation: sparkle3 10s ease-in-out infinite 4s;
}

@keyframes sparkle1 {
  0%,
  100% {
    opacity: 0;
    transform: scale(0.5);
  }
  50% {
    opacity: 1;
    transform: scale(1.2);
  }
}

@keyframes sparkle2 {
  0%,
  100% {
    opacity: 0;
    transform: scale(0.3) translateY(0px);
  }
  50% {
    opacity: 0.8;
    transform: scale(1) translateY(-5px);
  }
}

@keyframes sparkle3 {
  0%,
  100% {
    opacity: 0;
    transform: scale(0.4) rotate(0deg);
  }
  50% {
    opacity: 0.9;
    transform: scale(1.1) rotate(180deg);
  }
}
</style>
