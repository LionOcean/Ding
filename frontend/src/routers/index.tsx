import { lazy, Suspense } from "react";
import { RouteObject, Navigate, useRoutes } from "react-router-dom";

const lazyLoad = (Children: React.LazyExoticComponent<any>) => {
	return (
		<Suspense fallback={<>loading...</>}>
			<Children />
		</Suspense>
	);
};

export const routesConfig: RouteObject[] = [
	{
		path: "/",
		element: <Navigate to='/home' />,
	},
	{
		path: "/home",
		element: lazyLoad(lazy(() => import("../pages/Home"))),
	},
	{
		path: "/download",
		element: lazyLoad(lazy(() => import("../pages/Download"))),
	},
	{
		path: "/upload",
		element: lazyLoad(lazy(() => import("../pages/Upload"))),
	},
	{
		path: "*",
		element: <Navigate to='/home' />,
	},
];

const Router = () => {
	const routes = useRoutes(routesConfig);
	return routes;
};

export default Router;
