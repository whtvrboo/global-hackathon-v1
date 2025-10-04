<template>
  <div
    @click="$emit('click')"
    class="group relative bg-white rounded-xl border border-gray-200 overflow-hidden 
           hover:shadow-lg hover:border-gray-300 transition-all duration-200 cursor-pointer"
  >
    <div class="aspect-[2/3] relative overflow-hidden bg-gray-100">
      <img
        v-if="book.cover_url"
        :src="book.cover_url"
        :alt="book.title"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-4xl">
        📚
      </div>
      
      <!-- Status Badge -->
      <div
        v-if="book.status"
        class="absolute top-2 right-2 px-2 py-1 text-xs font-medium rounded-full bg-white/90 backdrop-blur"
        :class="statusColor(book.status)"
      >
        {{ statusLabel(book.status) }}
      </div>
    </div>

    <div class="p-4">
      <h3 class="font-semibold text-gray-900 line-clamp-2 mb-1">
        {{ book.title }}
      </h3>
      <p v-if="book.authors?.length" class="text-sm text-gray-600 line-clamp-1">
        {{ book.authors.join(', ') }}
      </p>
      
      <!-- Rating -->
      <div v-if="book.rating" class="flex items-center gap-1 mt-2">
        <span class="text-yellow-500">★</span>
        <span class="text-sm font-medium">{{ book.rating }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  book: {
    type: Object,
    required: true
  }
})

defineEmits(['click'])

const statusLabel = (status) => {
  const labels = {
    'want_to_read': 'Want to Read',
    'reading': 'Reading',
    'read': 'Read',
    'dnf': 'DNF'
  }
  return labels[status] || status
}

const statusColor = (status) => {
  const colors = {
    'want_to_read': 'text-blue-700',
    'reading': 'text-green-700',
    'read': 'text-purple-700',
    'dnf': 'text-gray-700'
  }
  return colors[status] || 'text-gray-700'
}
</script>

