<script lang="ts" setup>
import { onMounted, ref } from 'vue'

const text = ref<string>()
const items = ref<any[]>();

async function createJob() {
  let response = await fetch("/api/jobs", {
    method: "POST",
    body: JSON.stringify({
      priority: 0,
      script: text.value,
      type: "script",
      resourceRequest: "gpu"
    })
  });
  await fetchJobs();
}
function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString();
}

async function deleteJob(jobId: number) {
  await fetch(`/api/jobs/${jobId}`, {
    method: "DELETE"
  });
  await fetchJobs();
}

async function fetchJobs() {
  const response = await fetch("/api/jobs");
  items.value = await response.json();
}

onMounted(async () => {
  await fetchJobs();
})
</script>

<template>
  <div>
    <h1>Clai</h1>
    <div>
      <TextArea class="overflow-scroll max-h-30rem" autoResize cols="80" rows="40" v-model="text" />
    </div>
    <div>
      <Button label="Queue" @click="createJob" />
    </div>
  </div>
  <div class="ml-4 p-datatable p-component">
    <h1>Jobs</h1>
    <table class="p-datatable-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Created At</th>
          <th>Priority</th>
          <th>Resource Request</th>
          <th>Type</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ formatDate(item.createdAt) }}</td>
          <td>{{ item.priority }}</td>
          <td>{{ item.resourceRequest }}</td>
          <td>{{ item.type }}</td>
          <td><Button label="Delete" @click="deleteJob(item.id)" /></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
