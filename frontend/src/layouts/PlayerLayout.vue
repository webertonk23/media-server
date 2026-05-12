<template>
  <div class="player-layout min-h-screen bg-black flex items-center justify-center">
    <!-- Minimal layout for fullscreen player - no header/footer -->
    <!--<router-view />-->
    <slot />
  </div>
</template>
<script setup lang="ts">
</script>
<style scoped>
.player-layout {
  /* Pure black background for cinematic player experience */
  background-color: #000000;
  /* Ensure full viewport coverage */
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}
</style>
