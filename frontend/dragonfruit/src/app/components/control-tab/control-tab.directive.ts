import { Directive, ViewContainerRef } from '@angular/core';

@Directive({
    selector: '[cTabDirect]'
})
export class ControlTabDirective {
    constructor(public viewContainerRef: ViewContainerRef) { }
}
