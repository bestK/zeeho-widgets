<script setup>
import VehicleWidget from './components/VehicleWidget.vue';
</script>

<template>
    <div class="app">
        <VehicleWidget />
    </div>
</template>

<style>
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

.app {
    width: 100vw;
    height: 100vh;
    background: #fff; 
    padding: 4px;
}
</style>