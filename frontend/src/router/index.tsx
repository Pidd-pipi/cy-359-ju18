import { createBrowserRouter, Navigate } from 'react-router-dom'
import App from '../App'
import Login from '../pages/Login'
import Register from '../pages/Register'
import ActivityList from '../pages/ActivityList'
import ActivityDetail from '../pages/ActivityDetail'
import TeamList from '../pages/TeamList'
import TeamDetail from '../pages/TeamDetail'
import Checkin from '../pages/Checkin'
import LeaderboardPage from '../pages/LeaderboardPage'
import Mall from '../pages/Mall'
import Redemptions from '../pages/Redemptions'
import Favorites from '../pages/Favorites'
import Profile from '../pages/Profile'
import ActivityManage from '../pages/ActivityManage'
import ProductsManage from '../pages/ProductsManage'
import Users from '../pages/Users'
import AuditLogs from '../pages/AuditLogs'
import { getToken } from '../utils/auth'

function RequireAuth({ children }: { children: React.ReactNode }) {
  if (!getToken()) return <Navigate to="/login" replace />
  return <>{children}</>
}

export const router = createBrowserRouter([
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <ActivityList /> },
      { path: 'activities/:id', element: <ActivityDetail /> },
      { path: 'leaderboard/:id', element: <RequireAuth><LeaderboardPage /></RequireAuth> },
      { path: 'checkin/:id', element: <RequireAuth><Checkin /></RequireAuth> },
      { path: 'teams', element: <RequireAuth><TeamList /></RequireAuth> },
      { path: 'teams/:id', element: <RequireAuth><TeamDetail /></RequireAuth> },
      { path: 'mall', element: <RequireAuth><Mall /></RequireAuth> },
      { path: 'redemptions', element: <RequireAuth><Redemptions /></RequireAuth> },
      { path: 'favorites', element: <RequireAuth><Favorites /></RequireAuth> },
      { path: 'profile', element: <RequireAuth><Profile /></RequireAuth> },
      { path: 'admin/activities', element: <RequireAuth><ActivityManage /></RequireAuth> },
      { path: 'admin/products', element: <RequireAuth><ProductsManage /></RequireAuth> },
      { path: 'admin/users', element: <RequireAuth><Users /></RequireAuth> },
      { path: 'admin/audit', element: <RequireAuth><AuditLogs /></RequireAuth> },
    ],
  },
  { path: '*', element: <Navigate to="/" replace /> },
])
