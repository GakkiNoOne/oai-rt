import { defineConfig } from 'umi';

export default defineConfig({
  npmClient: 'pnpm',
  history: {
    type: 'hash',  // 使用 hash 路由
  },
  base: '/',
  publicPath: '/',
  routes: [
    { 
      path: '/login', 
      component: '@/pages/Login',
      layout: false,  // 登录页面不使用布局
    },
    {
      path: '/',
      component: '@/layouts/index',
      layout: false,  // 禁用 umi 内置 layout，使用自定义 layout
      routes: [
        { path: '/', redirect: '/rts' },
        { path: '/rts', component: '@/pages/RTs' },
        { path: '/configs', component: '@/pages/Configs' },
        { path: '/api-docs', component: '@/pages/ApiDocs' },
      ],
    },
  ],
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
});
