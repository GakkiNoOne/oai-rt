import { defineConfig } from 'umi';

// 使用绝对路径 / 避免 UMI 的开发环境限制
// 前端使用 hash 路由，资源都是相对路径加载，不会有跨域问题
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
