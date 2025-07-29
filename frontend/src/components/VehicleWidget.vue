<template>
    <div class="widget-container" style="--wails-draggable: drag">
        <!-- 标题栏 -->
        <div class="widget-header">
            <div class="title">ZEEHO ({{ vehicleDataList.length }})</div>
            <div class="actions">
                <!-- <button class="action-btn no-drag" @click="showConfirm('minimize')">
          _
        </button> -->
                <button class="action-btn no-drag" @click="showConfirm('exit')">×</button>
            </div>
        </div>

        <div class="widget-body" style="--wails-draggable: no-drag">
            <div v-if="loading" class="loading">
                <div class="spinner"></div>
            </div>

            <div v-else-if="error" class="error">
                <p>{{ error }}</p>
                <button @click="fetchData" class="retry-btn">重试</button>
            </div>

            <div v-else-if="vehicleDataList.length > 0" class="widget-content">
                <!-- 车辆滚动容器 -->
                <div class="vehicles-scroll-container">
                    <div class="vehicles-list">
                        <div
                            v-for="(vehicle, index) in vehicleDataList"
                            :key="vehicle.vehicleId || index"
                            class="vehicle-card"
                        >
                            <!-- 车辆名称 -->
                            <div class="vehicle-name">
                                {{ vehicle.vehicleName || '未知车辆' }}
                            </div>

                            <!-- 主要内容区域 -->
                            <div class="main-content">
                                <!-- 左侧信息 -->
                                <div class="info-section">
                                    <div class="stats">
                                        <span class="range">{{ vehicle.hmiRidableMile }}km</span>
                                        <span class="battery">{{ vehicle.bmssoc }}%</span>
                                    </div>
                                    <div class="battery-bar">
                                        <div
                                            class="battery-fill"
                                            :style="{ width: parseInt(vehicle.bmssoc || 0) + '%' }"
                                        ></div>
                                        <!-- 充电指示器 -->
                                        <div
                                            v-if="vehicle.chargeState === '1'"
                                            class="charging-indicator"
                                            :style="{ left: parseInt(vehicle.bmssoc || 0) + '%' }"
                                        >
                                            ⚡
                                        </div>
                                    </div>
                                </div>

                                <!-- 右侧车辆图片 -->
                                <div class="vehicle-section">
                                    <img
                                        v-if="vehicle.vehiclePicUrl"
                                        :src="vehicle.vehiclePicUrl"
                                        alt="Vehicle"
                                        class="vehicle-image"
                                        @error="handleImageError"
                                    />
                                    <div v-else class="vehicle-placeholder">🛵</div>
                                </div>
                            </div>

                            <!-- 位置信息 -->
                            <div class="vehicle-location">
                                <div class="location-time">📍 {{ vehicle.location?.locationTime || '位置信息' }}</div>
                                <div v-if="vehicle.location?.address" class="location-address">
                                    {{ vehicle.location?.address }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 底部控制按钮 -->
        <div class="widget-footer">
            <div class="controls">
                <button class="control-btn settings-btn" @click="openConfigModal" title="设置">⚙️</button>
                <button class="control-btn widget-btn" @click="startWidget">🧩</button>
                <button class="control-btn refresh-btn" @click="refreshWidget">⟳</button>
                <!-- <div class="vehicle-count">共 {{ vehicleDataList.length }} 台车辆</div> -->
                <div class="copyright">Power by KK</div>
            </div>
        </div>
        <!-- 配置模态框 -->
        <ConfigModal :show="showConfigModal" @close="closeConfigModal" @saved="onConfigSaved" />

        <!-- 确认对话框 -->
        <ConfirmDialog
            :show="confirmDialog.show"
            :title="confirmDialog.title"
            :message="confirmDialog.message"
            :confirm-text="confirmDialog.confirmText"
            @confirm="handleConfirm"
            @cancel="hideConfirm"
        />
    </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue';
import { GetConfig, Quit, ScheduleRefresh, StartWidget, VehicleHomePage } from '../../wailsjs/go/main/App';
import { EventsOff, EventsOn, WindowMinimise } from '../../wailsjs/runtime/runtime';
import ConfigModal from './ConfigModal.vue';
import ConfirmDialog from './ConfirmDialog.vue';

const scheduledId = ref(null);
const _config = ref();
const vehicleDataList = ref([]);
const loading = ref(true);
const error = ref(null);
const showConfigModal = ref(false);

// 确认对话框状态
const confirmDialog = ref({
    show: false,
    type: '',
    title: '',
    message: '',
    confirmText: '',
});

// 显示确认对话框
const showConfirm = type => {
    confirmDialog.value = {
        show: true,
        type,
        title: type === 'minimize' ? '最小化' : '退出程序',
        message: type === 'minimize' ? '确定要最小化小部件吗？' : '确定要退出程序吗？',
        confirmText: type === 'minimize' ? '最小化' : '退出',
    };
};

// 隐藏确认对话框
const hideConfirm = () => {
    confirmDialog.value.show = false;
};

// 处理确认操作
const handleConfirm = async () => {
    if (confirmDialog.value.type === 'minimize') {
        await WindowMinimise();
    } else {
        await Quit(); // 使用后端提供的 Quit 方法
    }
    hideConfirm();
};

// 保留其他现有方法
const fetchData = async () => {
    loading.value = true;
    error.value = null;
    try {
        vehicleDataList.value = await VehicleHomePage();
    } catch (err) {
        error.value = err.message || '获取数据失败';
    } finally {
        loading.value = false;
    }
};

const handleImageError = e => {
    e.target.style.display = 'none';
    e.target.nextElementSibling.style.display = 'flex';
};

const openConfigModal = () => {
    showConfigModal.value = true;
};

const closeConfigModal = () => {
    showConfigModal.value = false;
};

const onConfigSaved = () => {
    fetchData();
};

const startWidget = async () => {
    try {
        await StartWidget();
    } catch (err) {
        console.error('启动小部件失败:', err);
    }
};

const refreshWidget = async () => {
    await fetchData();
};

const initWidgets = async () => {
    const config = await GetConfig();

    if (config?.token) {
        _config.value = config;
        console.log('initWidgets', _config.value);
        fetchData();
    } else {
        showConfigModal.value = true;
    }
};

onMounted(async () => {
    // Set up event listeners first, before any initialization
    EventsOn('configUpdate', async function (data) {
        console.log('configUpdate', data);
        await initWidgets();
        await ScheduleRefresh();
    });

    EventsOn('dataRefreshed', function (data) {
        console.log('dataRefreshed', data);
        vehicleDataList.value = data;
    });

    EventsOn('refreshError', function (data) {
        console.log('refreshError', data);
    });

    // Now initialize widgets after listeners are set up
    await initWidgets();
});

onUnmounted(() => {
    EventsOff('configUpdate');
    EventsOff('dataRefreshed');
    EventsOff('refreshError');
});
</script>

<style scoped>
.widget-container {
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255);
    backdrop-filter: blur(40px) saturate(180%);
    -webkit-backdrop-filter: blur(40px) saturate(180%);
    border: 1px solid rgba(255, 255, 255, 0.18);
    border-radius: 20px;
    padding: 12px;
    box-sizing: border-box;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    position: relative;
    box-shadow:
        0 8px 32px rgba(0, 0, 0, 0.12),
        inset 0 1px 0 rgba(255, 255, 255, 0.3);
    overflow: hidden;

    /* From https://css.glass */
    background: rgba(39, 35, 35, 0.56);
    border-radius: 16px;
    box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(39, 35, 35, 0.65);
}

.widget-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    cursor: move;
    user-select: none;
}

.widget-header .title {
    font-size: 14px;
    font-weight: 600;
    color: #eee;
    opacity: 1;
}

.widget-header .actions {
    display: flex;
    gap: 8px;
}

.widget-header .action-btn {
    background: none;
    border: none;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    padding: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
}

.widget-header .action-btn:hover {
    background-color: #007aff;
    color: #fff;
}

.widget-body {
    position: relative;
    z-index: 3;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 75%;
    width: 100%;
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
    width: 100%;
    display: flex;
    flex-direction: column;
    position: relative;
    z-index: 3;
    overflow: hidden;
}

.vehicles-scroll-container {
    flex: 1;
    overflow-x: auto;
    overflow-y: hidden;
    padding: 8px 0;
}

.vehicles-list {
    display: flex;
    gap: 16px;
    padding: 0 8px;
    min-height: 100%;
}
.vehicle-card {
    position: relative;
    flex-shrink: 0;
    width: 280px;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    background-color: rgba(0, 0, 0, 0.5);
    border-radius: 12px;
    box-shadow: 0 0 30px rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(20px) saturate(180%);
    color: #f0f0f0;
    overflow: hidden;
}

/* 光斑1 */
.vehicle-card::before,
.vehicle-card::after {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle at center, rgba(255, 255, 255, 0.15) 0%, transparent 60%);
    pointer-events: none;
    filter: blur(40px);
    border-radius: 50%;
}

/* 第一个光斑：慢速移动 + 呼吸感 */
.vehicle-card::before {
    animation: glow1 16s ease-in-out infinite;
    opacity: 0.4;
}

/* 第二个光斑：稍微快点，路径不同 */
.vehicle-card::after {
    animation: glow2 12s ease-in-out infinite;
    opacity: 0.3;
}

@keyframes glow1 {
    0% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0.3;
    }
    50% {
        transform: translate(20%, 20%) scale(1.05);
        opacity: 0.5;
    }
    100% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0.3;
    }
}

@keyframes glow2 {
    0% {
        transform: translate(10%, -30%) scale(1);
        opacity: 0.25;
    }
    50% {
        transform: translate(-10%, 30%) scale(1.1);
        opacity: 0.4;
    }
    100% {
        transform: translate(10%, -30%) scale(1);
        opacity: 0.25;
    }
}
.vehicle-name {
    font-size: 14px;
    font-weight: 600;
    color: #333;
    text-align: center;
    padding-bottom: 8px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.vehicle-location {
    font-size: 10px;
    color: #666;
    text-align: center;
    margin-top: auto;
}

.vehicle-location .location-time {
    margin-bottom: 2px;
}

.vehicle-location .location-address {
    font-size: 9px;
    opacity: 0.8;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
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
    gap: 8px;
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
    overflow: visible;
    margin-bottom: 12px;
    position: relative;
}

.battery-fill {
    height: 100%;
    background: linear-gradient(90deg, #00d4ff 0%, #0099cc 100%);
    border-radius: 2px;
    transition: width 0.3s ease;
}

.charging-indicator {
    position: absolute;
    top: -8px;
    transform: translateX(-50%);
    font-size: 16px;
    animation: breathe 2s ease-in-out infinite;
    filter: drop-shadow(0 0 4px rgba(255, 193, 7, 0.8));
    z-index: 10;
}

@keyframes breathe {
    0%,
    100% {
        opacity: 0.6;
        transform: scale(1);
    }
    50% {
        opacity: 1;
        transform: scale(1.2);
    }
}

.vehicle-section {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
}

.vehicle-image {
    max-width: 120px;
    max-height: 80px;
    object-fit: contain;
}

.vehicle-placeholder {
    font-size: 32px;
    opacity: 0.6;
}

.widget-footer {
    display: flex;
    justify-content: space-between;
    height: 15%;
    width: 100%;
}

.widget-footer .controls {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-top: 8px;
    width: 100%;
}

.control-btn {
    width: 32px;
    height: 32px;
    background: rgba(5, 5, 5, 0.56);
    color: #fff;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    cursor: pointer;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
}

.refresh-btn:hover {
    color: #333;
}

.control-btn:hover {
    background: #f0f0f0;
}

.vehicle-count {
    flex: 1;
    font-size: 12px;
    color: #666;
    text-align: right;
    display: flex;
    align-items: center;
    justify-content: flex-end;
}

.copyright {
    flex: 1;
    font-size: 12px;
    color: #666;
    text-align: right;
    display: flex;
    align-items: center;
    justify-content: flex-end;
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

/* 添加滚动条样式 */
.vehicles-scroll-container::-webkit-scrollbar {
    height: 6px;
}

.vehicles-scroll-container::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 3px;
}

.vehicles-scroll-container::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.3);
    border-radius: 3px;
}

.vehicles-scroll-container::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.5);
}
</style>
