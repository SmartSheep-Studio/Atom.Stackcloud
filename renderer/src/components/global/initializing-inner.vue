<template>
  <div class="wrapper">
    <div style="min-width: 180px">
      <v-card variant="flat" color="transparent">
        <v-card-text>
          <div>Initializing...</div>
        </v-card-text>
      </v-card>

      <v-progress-linear indeterminate />
    </div>
  </div>
</template>

<style scoped>
.wrapper {
  height: 100%;
  width: 100%;

  text-align: center;
  display: flex;
  place-items: center;
  justify-content: center;

  transition: all 300ms;
}
</style>
