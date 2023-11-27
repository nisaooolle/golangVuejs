<template>
    <h1>{{ title }}</h1>

<div class="product-card">
    <img src="../assets/images.jpg" class="product-image" alt="">

    <div class="product-detail">
        <h2 class="product-name">{{ title }}</h2>
        <h2  class="product-description">{{ description }}</h2>
        <h2  class="product-price">{{ price }}</h2>
    </div>
</div>
</template>

<script>
export default {
    data() {
        return {
            title: "E-commerce"
        }
    }
}
</script>