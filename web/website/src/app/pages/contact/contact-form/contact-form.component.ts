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

import { Component, OnDestroy, ViewChild, ViewEncapsulation } from "@angular/core";
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
import { ApiError, ApiErrorCode, ApiErrorMessage } from "../../../../../../pkg/api/errors";

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
export class ContactFormComponent implements OnDestroy {
    form: FormGroup;

    @ViewChild('formDirective')
    formDirective!: FormGroupDirective;

    nameMatcher = new ContactFormErrorStateMatcher();
    emailMatcher = new ContactFormErrorStateMatcher();
    msgMatcher = new ContactFormErrorStateMatcher();

    sending = false;

    private _destroyed = false;

    constructor(private _contactFormService: ContactFormService, private _snackBar: MatSnackBar) {
        this.form = new FormGroup(
            {
                name: new FormControl('', [Validators.required, Validators.maxLength(100), createWhitespaceValidator()]),
                email: new FormControl('', [Validators.required, Validators.email, Validators.maxLength(500)]),
                msg: new FormControl('', [Validators.required, Validators.maxLength(1000), createWhitespaceValidator()])
            }
        );
    }

    ngOnDestroy(): void {
        this._destroyed = true;
    }

    async submit(): Promise<void> {
        if (this.sending) {
            return;
        }
        if (this.form.invalid) {
            this.form.updateValueAndValidity();
            return;
        }

        this.sending = true;
        let v = this.form.value;
        let msg: Message = {
            name: v.name,
            email: v.email,
            message: v.msg
        };

        try {
            await this._contactFormService.send(msg);
        } catch (e) {
            if (!this._destroyed) {
                this.sending = false;
                let msg = '';
                if (e instanceof ApiError) {
                    switch (e.code) {
                        case ApiErrorCode.BAD_REQUEST:
                            msg = 'Error: ' + (e.message ? e.message : ApiErrorMessage.INVALID_OPERATION);
                            break;
                        case ApiErrorCode.PERMISSION_DENIED:
                            msg = 'Error: ' + ApiErrorMessage.PERMISSION_DENIED;
                            break;
                        case ApiErrorCode.INVALID_DATA:
                            msg = 'Error: ' + (e.message ? e.message : ApiErrorMessage.INVALID_OPERATION);
                            break;
                        default:
                            msg = 'An error occurred while sending a message';
                    }
                }

                this._snackBar.open(msg, 'X', { duration: SNACK_BAR_DURATION });
            }
            return;
        }

        if (this._destroyed) {
            return;
        }

        this.sending = false;
        this._snackBar.open('The message has been sent', 'X', { duration: SNACK_BAR_DURATION });
        this.resetForm();
    }

    private resetForm(): void {
        this.formDirective.resetForm({
            name: '',
            email: '',
            msg: ''
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
