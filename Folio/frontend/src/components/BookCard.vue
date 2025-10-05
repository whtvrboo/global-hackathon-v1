<template>
  <div @click="$emit('click')" class="group relative bg-dark-900/50 backdrop-blur-sm rounded-2xl border border-dark-800 overflow-hidden 
           hover:bg-dark-800/50 hover:border-dark-700 hover:shadow-2xl hover:shadow-black/20 
           transition-all duration-300 cursor-pointer transform hover:scale-[1.02]">
    <div class="aspect-[2/3] relative overflow-hidden bg-dark-800">
      <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
      <div v-else class="w-full h-full flex items-center justify-center text-4xl text-dark-400">
      </div>

      <!-- Status Badge -->
      <div v-if="book.status" class="absolute top-3 right-3 px-3 py-1 text-xs font-semibold rounded-full glass-strong"
        :class="statusColor(book.status)">
        {{ statusLabel(book.status) }}
      </div>

      <!-- Rating Overlay -->
      <div v-if="book.rating"
        class="absolute bottom-3 left-3 flex items-center gap-1 px-2 py-1 glass-strong rounded-full">
        <span class="text-accent-orange text-sm">★</span>
        <span class="text-white text-sm font-semibold">{{ book.rating }}</span>
      </div>
    </div>

    <div class="p-4">
      <h3 class="font-semibold text-white line-clamp-2 mb-2 text-sm leading-tight">
        {{ book.title }}
      </h3>
      <p v-if="book.authors?.length" class="text-xs text-dark-400 line-clamp-1">
        {{ book.authors.join(', ') }}
      </p>
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
    'want_to_read': 'text-accent-blue bg-accent-blue/20',
    'reading': 'text-accent-green bg-accent-green/20',
    'read': 'text-accent-purple bg-accent-purple/20',
    'dnf': 'text-dark-400 bg-dark-600/20'
  }
  return colors[status] || 'text-dark-400 bg-dark-600/20'
}
</script>
