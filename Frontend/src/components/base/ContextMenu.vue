<template>
    <div v-if="visible" class="context-menu" :style="{ top: `${y}px`, left: `${x}px` }">
        <ul>
        <li v-for="item in items" :key="item" @click="handleClick(item)">
            {{ item }}
        </li>
        </ul>
    </div>
</template>

<script>
export default {
    data() {
        return {
        visible: false,
        x: 0,
        y: 0,
        items: [],
        obj: null,  // 当前触发菜单的对象
        };
    },
    methods: {
        show(event, x, y, items, obj) {
            console.log('show');
            this.x = x;
            this.y = y;
            this.items = items;
            this.visible = true;
            this.obj = obj;
            // 使用 setTimeout 延迟添加点击监听器（防止使用点击触发菜单时，菜单被直接隐藏）
            setTimeout(() => {
                document.addEventListener('click', this.hide);
            }, 0);
        },
        hide() {
            this.visible = false;
            document.removeEventListener('click', this.hide);
        },
        // 选中某个选项
        handleClick(item) {
            this.$emit('select-item', item, this.obj);
            this.hide();
        },
    },
};
</script>

<style scoped>
.context-menu {
    position: absolute;
    background-color: white;
    border: 1px solid #ccc;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
    z-index: 1000;
}
.context-menu ul {
    list-style: none;
    margin: 0;
    padding: 0;
}
.context-menu li {
    padding: 8px 12px;
    cursor: pointer;
}
.context-menu li:hover {
    background-color: #f5f5f5;
}
</style>