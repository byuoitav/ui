import { Component, OnInit, ViewChild, forwardRef, Input, ComponentFactoryResolver, ViewContainerRef, AfterViewInit } from '@angular/core';
import { ControlTabDirective } from './control-tab.directive';
import { IControlTab } from './icontrol-tab';
import { ControlGroup, CONTROL_TAB, AUDIO_TAB, PRESENT_TAB, HELP_TAB } from '../../objects/control';
import { SingleDisplayComponent } from '../single-display/single-display.component';
import { MultiDisplayComponent } from '../multi-display/multi-display.component';
import { AudioComponent } from '../audio/audio.component';
import { PresentComponent } from '../present/present.component';
import { HelpComponent } from '../help/help.component';

@Component({
    selector: 'control-tab',
    templateUrl: './control-tab.component.html',
    styleUrls: ['./control-tab.component.scss']
})
export class ControlTabComponent implements OnInit, AfterViewInit {
    @ViewChild(ControlTabDirective, {static: true}) direct: ControlTabDirective;
    @Input() controlGroup: ControlGroup;
    @Input() tab: string;

    constructor(
        private resolver: ComponentFactoryResolver
    ) {}

    ngOnInit() {
        
    }

    ngAfterViewInit() {
        console.log('initializing the control tab component', this.controlGroup, this.tab);
        this.setUpComponent();
    }

    ngOnChanges() {
        this.setUpComponent();
    }

    private setUpComponent() {
        if (this.controlGroup && this.tab) {
            let comp: any;

            switch (this.tab) {
                case CONTROL_TAB:
                    if (this.controlGroup.displays.length < 2) {
                        comp = SingleDisplayComponent;
                    } else {
                        comp = MultiDisplayComponent;
                    }
                    break;
                case AUDIO_TAB:
                    comp = AudioComponent;
                    break;
                case PRESENT_TAB:
                    comp = PresentComponent;
                    break;
                case HELP_TAB:
                    comp = HelpComponent;
                    break;
                default:
                    break;
            }

            const factory = this.resolver.resolveComponentFactory(comp);
            const viewContainerRef = this.direct.viewContainerRef;
            viewContainerRef.clear();

            const componentRef = viewContainerRef.createComponent(factory);

            (<IControlTab>(componentRef.instance)).cg = this.controlGroup;
        }
    }
}
