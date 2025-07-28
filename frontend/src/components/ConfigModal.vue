<template>
    <div v-if="show" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
            <div class="modal-header">
                <h2>配置设置</h2>
                <button @click="closeModal" class="close-btn">&times;</button>
            </div>

            <div class="modal-body">
                <div class="form-group">
                    <label for="token">Token:</label>
                    <input
                        id="token"
                        v-model="formData.token"
                        type="text"
                        placeholder="请输入您的Token"
                        class="form-input"
                        :disabled="loading"
                    />
                </div> 
                <div class="form-group">
                    <label for="updateInterval">更新间隔（分钟）:</label>
                    <input
                        id="updateInterval"
                        v-model.number="formData.updateInterval"
                        type="number"
                        min="1"
                        max="60"
                        class="form-input"
                        :disabled="loading"
                    />
                    <small class="form-hint">建议设置在1-60分钟之间</small>
                </div>

                <div v-if="error" class="error-message">
                    {{ error }}
                </div>

                <div v-if="success" class="success-message">配置保存成功！</div>
            </div>

            <div class="modal-footer">
                <button @click="closeModal" class="btn btn-secondary" :disabled="loading">取消</button>
                <button @click="saveConfig" class="btn btn-primary" :disabled="loading || !canSave">
                    <span v-if="loading">验证中...</span>
                    <span v-else>保存配置</span>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue';
import { GetConfig, ValidateAndSaveConfig } from '../../wailsjs/go/main/App';

const props = defineProps({
    show: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(['close', 'saved']);

const formData = ref({
    token: '',
    vehicleId: '',
    updateInterval: 5, // 默认5分钟
});

const loading = ref(false);
const error = ref('');
const success = ref(false);

const canSave = computed(() => {
    const interval = parseInt(formData.value.updateInterval);
    return (
        formData.value.token.trim() !== '' &&
        formData.value.vehicleId.trim() !== '' &&
        !isNaN(interval) &&
        interval >= 1 &&
        interval <= 60
    );
});

const closeModal = () => {
    if (!loading.value) {
        emit('close');
    }
};

const saveConfig = async () => {
    if (!canSave.value) return;

    loading.value = true;
    error.value = '';
    success.value = false;

    try {
        // 验证更新间隔
        const interval = parseInt(formData.value.updateInterval);
        if (isNaN(interval) || interval < 1 || interval > 60) {
            throw new Error('更新间隔必须在1-60分钟之间');
        }

        await ValidateAndSaveConfig(formData.value.token.trim(), formData.value.vehicleId.trim(), interval);
        success.value = true;

        setTimeout(() => {
            emit('saved');
            closeModal();
        }, 1500);
    } catch (err) {
        error.value = err.message || '保存配置失败';
    } finally {
        loading.value = false;
    }
};

const loadCurrentConfig = async () => {
    try {
        const config = await GetConfig();
        if (config) {
            formData.value.token = config.token || '';
            formData.value.vehicleId = config.vehicleId || '';
        }
    } catch (err) {
        console.error('加载配置失败:', err);
    }
};

watch(
    () => props.show,
    newShow => {
        if (newShow) {
            error.value = '';
            success.value = false;
            loadCurrentConfig();
        }
    },
);
</script>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.modal-content {
    background: #f8f8f8;
    border-radius: 8px;
    width: 280px;
    height: auto;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid #ddd;
}

.modal-header h2 {
    margin: 0;
    color: #333;
    font-size: 14px;
    font-weight: 600;
}

.close-btn {
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
}

.close-btn:hover {
    color: #333;
}

.modal-body {
    padding: 12px;
}

.form-group {
    margin-bottom: 8px;
}

.form-group label {
    display: block;
    margin-bottom: 4px;
    color: #333;
    font-size: 11px;
    font-weight: 500;
}

.form-input {
    width: 100%;
    padding: 6px 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: white;
    color: #333;
    font-size: 11px;
    transition: border-color 0.3s;
    box-sizing: border-box;
}

.form-input:focus {
    outline: none;
    border-color: #007aff;
}

.form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.error-message {
    background: #ff3b30;
    color: white;
    padding: 6px 8px;
    border-radius: 4px;
    margin-top: 8px;
    font-size: 10px;
}

.success-message {
    background: #34c759;
    color: white;
    padding: 6px 8px;
    border-radius: 4px;
    margin-top: 8px;
    font-size: 10px;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 6px;
    padding: 8px 12px;
    border-top: 1px solid #ddd;
}

.btn {
    padding: 4px 8px;
    border: none;
    border-radius: 4px;
    font-size: 10px;
    cursor: pointer;
    transition: background-color 0.3s;
}

.btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-secondary {
    background: #f0f0f0;
    color: #333;
}

.btn-secondary:hover:not(:disabled) {
    background: #e0e0e0;
}

.btn-primary {
    background: #007aff;
    color: white;
}

.btn-primary:hover:not(:disabled) {
    background: #0056cc;
}

.form-hint {
    display: block;
    color: #666;
    font-size: 10px;
    margin-top: 2px;
}

input[type='number'] {
    -moz-appearance: textfield;
}

input[type='number']::-webkit-outer-spin-button,
input[type='number']::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
</style>
