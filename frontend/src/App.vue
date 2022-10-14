<script setup lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
import { ref } from 'vue';

import Layout from './components/layout/Layout.vue'

interface Email {
  id: number;
  from: string;
  to: string;
  subject: string;
  content: string;
  date: string;
  read: boolean;
}

fetch('http://localhost:8000/search')
  .then(response => response.json())
  .then(json => {
    emails.value = json.hits.hits.map((hit: any) => {
      return {
        id: hit._id,
        from: hit._source.From,
        to: hit._source.To,
        subject: hit._source.Subject,
        content: hit._source.Content,
        date: hit._source.Date
      }
    });
    count.value = json.hits.total.value;
    
    
  })

const selectedEmail = ref<Email | null>(null)

const emails = ref<Email[]>([])
const query = ref('')
const count = ref(0)
const offset = ref(0)
const limit = ref(20)

const selectEmail = (email: Email) => {
  selectedEmail.value = email
}

const setQuery = (e: Event) => {
  query.value = (e.target as HTMLInputElement).value
}

const submitSearch = (e: Event) => {
  offset.value = 0
  search(e)
}

const search = (e: Event) => {
  e.preventDefault()
  
  fetch(`http://localhost:8000/search?q=${query.value}&limit=${limit.value}&offset=${offset.value}`)
    .then(response => response.json())
    .then(json => {
      emails.value = json.hits.hits.map((hit: any) => {
        return {
          id: hit._id,
          from: hit._source.From,
          to: hit._source.To,
          subject: hit._source.Subject,
          content: hit._source.Content,
          date: hit._source.Date
        }
      });
      count.value = json.hits.total.value;
    })
}

const next = () => {
  offset.value += limit.value
  search(new Event('submit'))
}

const prev = () => {
  offset.value -= limit.value
  search(new Event('submit'))
}
</script>

<template>
  
  <Layout>
    <div class="min-h-screen">
      <label for="default-search" class="mb-2 text-sm font-medium text-gray-900 sr-only">Search</label>
      <form @submit="(e)=> submitSearch(e)">
        <div class="relative">
            <input type="search" :value="query" @change="(e)=> setQuery(e)" name="search" id="default-search" class="block p-4 pl-10 w-full text-sm text-gray-900 bg-gray-50 border border-gray-300 focus:ring-blue-500 focus:border-blue-500" placeholder="Search">
            <button type="submit" class="text-white absolute right-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 "><svg aria-hidden="true" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg></button>
        </div>
      </form>
      <div class="flex mb-4">
        <div class="w-1/2">
          <table class="border-collapse table-fixed break-words w-full bg-white text-sm shadow-sm">
            <thead class="bg-gray-400 border-gray-900">
              <tr>
                <th class="border border-slate-300 p-4 text-slate-900 text-left">Subject</th>
                <th class="border border-slate-300 p-4 text-slate-900 text-left">From</th>
                <th class="border border-slate-300 p-4 text-slate-900 text-left">To</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="email in emails" :key="email.id" @click="()=>selectEmail(email)" class="hover:bg-slate-100 cursor-pointer">
                <td class="border border-slate-300  text-slate-900 text-left p-4">
                  <span class="font-semibold ">{{ email.subject ? email.subject : "[No Subject]" }}</span>
                </td>
                <td class="border border-slate-300 text-slate-900  text-left p-4">
                  <span class="font-semibold">{{ email.from }}</span>
                </td>
                <td class="border border-slate-300  text-slate-900  text-left p-4">
                  <span class="font-semibold">{{ email.to }}</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="w-1/2">
          <div v-if="selectedEmail" class="p-10">
            <div class="mb-3">
              <span class="font-semibold ">{{selectedEmail.subject ? selectedEmail.subject: "[No Subject]"}}</span>
            </div>
            <div v-html="selectedEmail.content">
            </div>
          </div>
        </div>
      </div>
      <div class="m-5 flex justify-between">
        <button @click="prev" v-if="offset !== 0">&lt;</button>
        <span class="font-semibold">Showing {{offset + 1}}-{{offset + limit < count ? offset + limit : count}} of {{count}} emails</span>
        <button @click="next" v-if="offset + limit <= count">&gt;</button>
      </div>
    </div>
  </Layout>
</template>
