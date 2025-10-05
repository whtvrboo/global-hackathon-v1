<template>
    <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0"
        enter-to-class="opacity-100" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100"
        leave-to-class="opacity-0">
        <div v-if="show" @click="$emit('close')"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
            <div @click.stop
                class="bg-dark-900 rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-8">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-2xl font-bold text-white">
                        {{ isEditing ? 'Edit List' : 'Create New List' }}
                    </h2>
                    <button @click="$emit('close')" class="text-dark-400 hover:text-white transition-colors">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                </div>

                <form @submit.prevent="handleSubmit" class="space-y-6">
                    <!-- Live Preview -->
                    <div>
                        <label class="block text-sm font-medium text-white mb-2">
                            Live Preview
                        </label>
                        <div
                            class="aspect-video rounded-xl mb-4 flex items-center justify-center p-4 text-center bg-dark-800 border-2 border-dark-700 relative overflow-hidden">
                            <div v-if="form.header_image_url" class="absolute inset-0">
                                <img :src="form.header_image_url" alt="Header preview"
                                    class="w-full h-full object-cover" @error="onImageError" />
                                <div class="absolute inset-0 bg-black/40"></div>
                            </div>
                            <div class="relative">
                                <h3 class="text-2xl font-bold" :style="{ color: form.theme_color }">
                                    {{ form.name || 'My Awesome List' }}
                                </h3>
                                <p v-if="form.description" class="text-sm text-white/80 mt-1 line-clamp-2">
                                    {{ form.description }}
                                </p>
                            </div>
                        </div>
                    </div>


                    <!-- List Name -->
                    <div>
                        <label class="block text-sm font-medium text-white mb-2">
                            List Name *
                        </label>
                        <Input v-model="form.name" placeholder="e.g., My Favorite Fantasy Books" required />
                    </div>

                    <!-- Description -->
                    <div>
                        <label class="block text-sm font-medium text-white mb-2">
                            Description
                        </label>
                        <TextArea v-model="form.description" placeholder="Describe what this list is about..."
                            :rows="3" />
                    </div>

                    <!-- Header Image URL -->
                    <div>
                        <label class="block text-sm font-medium text-white mb-2">
                            Header Image URL
                        </label>
                        <Input v-model="form.header_image_url" placeholder="https://example.com/image.jpg" />
                        <p class="text-xs text-dark-400 mt-1">Optional: Add a beautiful header image for your list</p>
                    </div>

                    <!-- Theme Color -->
                    <div>
                        <label class="block text-sm font-medium text-white mb-2">
                            Theme Color
                        </label>
                        <div class="flex items-center gap-4">
                            <input v-model="form.theme_color" type="color"
                                class="w-12 h-12 rounded-lg border-2 border-dark-700 cursor-pointer" />
                            <Input v-model="form.theme_color" placeholder="#6366f1" class="flex-1" />
                        </div>
                        <p class="text-xs text-dark-400 mt-1">Choose a color that represents your list's theme</p>
                    </div>

                    <!-- Visibility -->
                    <div class="flex items-center gap-3">
                        <input v-model="form.is_public" type="checkbox" id="is_public"
                            class="w-4 h-4 text-accent-blue border-dark-700 rounded focus:ring-accent-blue bg-dark-800" />
                        <label for="is_public" class="text-sm text-white">
                            Make this list public (visible to other users)
                        </label>
                    </div>

                    <!-- Error -->
                    <div v-if="error"
                        class="p-4 bg-accent-red/10 border border-accent-red/20 rounded-lg text-accent-red text-sm">
                        {{ error }}
                    </div>

                    <!-- Actions -->
                    <div class="flex gap-3 justify-end pt-4 border-t border-dark-800">
                        <SecondaryButton @click="$emit('close')" type="button">
                            Cancel
                        </SecondaryButton>
                        <PrimaryButton type="submit" :loading="loading">
                            {{ isEditing ? 'Update List' : 'Create List' }}
                        </PrimaryButton>
                    </div>
                </form>
            </div>
        </div>
    </transition>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import axios from 'axios'
import PrimaryButton from './ui/PrimaryButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'
import Input from './ui/Input.vue'
import TextArea from './ui/TextArea.vue'
import { useToastStore } from '../stores/toast'

const props = defineProps({
    show: Boolean,
    list: Object, // For editing existing list
    book: Object,  // For adding book to new list
    prefilledData: Object // For pre-filling from thread
})

const emit = defineEmits(['close', 'success'])

const form = reactive({
    name: '',
    description: '',
    is_public: true,
    header_image_url: '',
    theme_color: '#6366f1'
})

const loading = ref(false)
const error = ref(null)
const isEditing = ref(false)
const toastStore = useToastStore()

watch(() => props.show, (newShow) => {
    if (newShow) {
        if (props.list) {
            // Editing existing list
            isEditing.value = true
            form.name = props.list.name || ''
            form.description = props.list.description || ''
            form.is_public = props.list.is_public !== false
            form.header_image_url = props.list.header_image_url || ''
            form.theme_color = props.list.theme_color || '#6366f1'
        } else if (props.prefilledData) {
            // Creating new list with pre-filled data from thread
            isEditing.value = false
            form.name = props.prefilledData.title || ''
            form.description = props.prefilledData.description || ''
            form.is_public = true
            form.header_image_url = ''
            form.theme_color = '#6366f1'
        } else {
            // Creating new list
            isEditing.value = false
            form.name = ''
            form.description = ''
            form.is_public = true
            form.header_image_url = ''
            form.theme_color = '#6366f1'
        }
        error.value = null
    }
})

const onImageError = (event) => {
    event.target.style.display = 'none'; // Hide broken image icon
}

const handleSubmit = async () => {
    if (!form.name.trim()) {
        error.value = 'List name is required'
        return
    }

    loading.value = true
    error.value = null

    try {
        if (isEditing.value) {
            // Update existing list
            await axios.put(`/api/lists/${props.list.id}`, {
                name: form.name,
                description: form.description,
                is_public: form.is_public,
                header_image_url: form.header_image_url || null,
                theme_color: form.theme_color
            })

            toastStore.success(`List "${form.name}" updated successfully!`)
        } else {
            // Create new list
            const response = await axios.post('/api/lists', {
                name: form.name,
                description: form.description,
                is_public: form.is_public,
                header_image_url: form.header_image_url || null,
                theme_color: form.theme_color
            })

            // If a book was provided, add it to the new list
            if (props.book?.id) {
                await axios.post(`/api/lists/${response.data.id}/items`, {
                    book_id: props.book.id
                })
            }

            // If pre-filled books were provided, add them to the new list
            if (props.prefilledData?.books && props.prefilledData.books.length > 0) {
                for (const bookData of props.prefilledData.books) {
                    await axios.post(`/api/lists/${response.data.id}/items`, {
                        book_id: bookData.id,
                        curator_note: bookData.annotations ?
                            bookData.annotations.map(a => a.content).join('\n\n') : ''
                    })
                }
            }

            toastStore.success(`List "${form.name}" created successfully!`)
        }

        emit('success')
        emit('close')
    } catch (err) {
        error.value = err.response?.data?.error || 'Failed to save list'
        toastStore.error('Failed to save list. Please try again.')
        console.error('Error saving list:', err)
    } finally {
        loading.value = false
    }
}
</script>
