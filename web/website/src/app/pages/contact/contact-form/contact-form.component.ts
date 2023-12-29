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

import { Component, ViewEncapsulation } from "@angular/core";
import {
    AbstractControl,
    FormControl,
    FormGroup,
    FormGroupDirective,
    NgForm,
    ValidationErrors,
    ValidatorFn,
    Validators,
    FormsModule,
    ReactiveFormsModule,
} from '@angular/forms';
import { ErrorStateMatcher } from '@angular/material/core';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar } from '@angular/material/snack-bar';

import { ContactFormService } from "./contact-form.service";
import { Message } from "./contact-form.api-models";

const SNACK_BAR_DURATION: number = 5000; // in milliseconds

@Component({
    selector: "pw-contact-form",
    standalone: true,
    templateUrl: './contact-form.component.html',
    styleUrls: ['./contact-form.component.css'],
    encapsulation: ViewEncapsulation.None,
    providers: [ContactFormService],
    imports: [FormsModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule, MatButtonModule]
})
export class ContactFormComponent {
    form: FormGroup;
    nameMatcher = new ContactFormErrorStateMatcher();
    emailMatcher = new ContactFormErrorStateMatcher();
    msgMatcher = new ContactFormErrorStateMatcher();

    message = new Message('', '', '');
    sending = false;

    constructor(private _contactFormService: ContactFormService, private _snackBar: MatSnackBar) {
        this.form = new FormGroup(
            {
                'name': new FormControl('', [Validators.required, Validators.maxLength(100), createWhitespaceValidator()]),
                'email': new FormControl('', [Validators.required, Validators.email, Validators.maxLength(500)]),
                'msg': new FormControl('', [Validators.required, Validators.maxLength(1000), createWhitespaceValidator()])
            }
        );
    }

    submit(): void {
        if (this.sending) {
            return;
        }
        if (this.form.invalid) {
            this.form.updateValueAndValidity();
            return;
        }

        this.sending = true;

        this._contactFormService.send(this.message).then(r => {
            this.sending = false;
            this._snackBar.open('The message has been sent.', 'X', { duration: SNACK_BAR_DURATION });
        }).catch(err => {
            console.error('[contact-form.ContactFormComponent.submit] send a message:', err);
            this.sending = false;
            this._snackBar.open('Send error.', 'X', { duration: SNACK_BAR_DURATION });
        });
    }
}

class ContactFormErrorStateMatcher implements ErrorStateMatcher {
    isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
        const isSubmitted = form && form.submitted;
        return !!(control && control.invalid && (control.dirty || control.touched || isSubmitted));
    }
}

function createWhitespaceValidator(): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
        const v = control.value;
        if (!v) {
            return null;
        }
        return /^[\s\u0085]+$/.test(v) ? { whitespace: true } : null;
    }
}
