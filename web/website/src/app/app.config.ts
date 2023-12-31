/**
 * @license
 * Copyright 2023 Alexey Lavrenchenko. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ApplicationConfig } from "@angular/core";
import { LocationStrategy, PathLocationStrategy } from "@angular/common";
import { provideHttpClient } from "@angular/common/http";
import { provideRouter, Routes, withInMemoryScrolling } from "@angular/router";
import { provideAnimations, provideNoopAnimations } from "@angular/platform-browser/animations";

import { HomeComponent } from "./pages/home";
import { InfoComponent } from "./pages/info";
import { AboutComponent } from "./pages/about";
import { ContactComponent } from "./pages/contact";
import { NotFoundComponent } from "./pages/not-found";

const appRoutes: Routes = [
    { title: "Alexey Lavrenchenko", path: "", pathMatch: "full", component: HomeComponent },
    { title: "Info", path: "info", pathMatch: "full", component: InfoComponent },
    { title: "About me", path: "about", pathMatch: "full", component: AboutComponent },
    { title: "Contact me", path: "contact", pathMatch: "full", component: ContactComponent },
    { title: "Not Found", path: "**", component: NotFoundComponent }
];

const prefersReducedMotion = typeof matchMedia === 'function' ? matchMedia('(prefers-reduced-motion)').matches : false;

export const appConfig: ApplicationConfig = {
    providers: [
        provideRouter(appRoutes, withInMemoryScrolling({
            scrollPositionRestoration: 'enabled',
            anchorScrolling: 'enabled'
        })),
        { provide: LocationStrategy, useClass: PathLocationStrategy },
        provideHttpClient(),
        prefersReducedMotion ? provideNoopAnimations() : provideAnimations(),
    ]
};
