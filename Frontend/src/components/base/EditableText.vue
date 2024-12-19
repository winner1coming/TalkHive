<template>
  <div 
    class="editable-text" 
    @mouseover="showEditIcon = true" 
    @mouseleave="showEditIcon = false"
    @click="startEditing"
  >
    <div v-if="!isEditing" class="text-display">
      <span>{{ text }}</span>
      <i v-if="showEditIcon" class="edit-icon" >✏️</i>
    </div>
    <div v-else class="text-edit">
      <input type="text" v-model="editableText" @blur="saveEdit" @keyup.enter="saveEdit" />
    </div>
  </div>
</template>

<script>
export default {
  props: {
    text: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      isEditing: false,
      showEditIcon: false,
      editableText: this.text
    };
  },
  methods: {
    startEditing() {
      this.isEditing = true;
      this.$nextTick(() => {
        this.$refs.editInput.focus();
      });
    },
    saveEdit() {
      this.isEditing = false;
      this.$emit('update-text', this.editableText);
    }
  },
  watch: {
    text(newText) {
      this.editableText = newText;
    }
  }
};
</script>

<style scoped>
.editable-text {
  position: relative;
  cursor: pointer;
}

.text-display {
  display: flex;
  align-items: flex-start;
}

.edit-icon {
  margin-left: 5px;
  cursor: pointer;
  font-size: 14px;
}

.text-edit input {
  width: 100%;
  padding: 5px;
  font-size: 14px;
}

span {
  text-align: left;
}
</style>