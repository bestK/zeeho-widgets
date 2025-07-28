<template>
    <div v-if="show" class="confirm-overlay" @click="onCancel">
        <div class="confirm-content" @click.stop>
            <div class="confirm-header">
                <h3>{{ title }}</h3>
            </div>
            <div class="confirm-body">
                {{ message }}
            </div>
            <div class="confirm-footer">
                <button class="btn btn-secondary" @click="onCancel">取消</button>
                <button class="btn btn-primary" @click="onConfirm">{{ confirmText }}</button>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps({
    show: {
        type: Boolean,
        default: false,
    },
    title: {
        type: String,
        default: '确认',
    },
    message: {
        type: String,
        required: true,
    },
    confirmText: {
        type: String,
        default: '确认',
    },
});

const emit = defineEmits(['cancel', 'confirm']);

const onCancel = () => emit('cancel');
const onConfirm = () => emit('confirm');
</script>

<style scoped>
.confirm-overlay {
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

.confirm-content {
    background: #f8f8f8;
    border-radius: 8px;
    width: 240px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.confirm-header {
    padding: 8px 12px;
    border-bottom: 1px solid #ddd;
}

.confirm-header h3 {
    margin: 0;
    color: #333;
    font-size: 13px;
    font-weight: 600;
}

.confirm-body {
    padding: 12px;
    color: #666;
    font-size: 12px;
    line-height: 1.4;
}

.confirm-footer {
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
    font-size: 11px;
    cursor: pointer;
    transition: background-color 0.3s;
}

.btn-secondary {
    background: #f0f0f0;
    color: #333;
}

.btn-secondary:hover {
    background: #e0e0e0;
}

.btn-primary {
    background: #007aff;
    color: white;
}

.btn-primary:hover {
    background: #0056cc;
}
</style>
