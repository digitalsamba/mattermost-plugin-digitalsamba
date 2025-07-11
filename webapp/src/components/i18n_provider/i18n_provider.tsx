import React from 'react';
import {IntlProvider} from 'react-intl';

interface Props {
    children: React.ReactNode;
}

export default function I18nProvider(props: Props) {
    return (
        <IntlProvider
            locale='en'
            messages={{}}
        >
            {props.children}
        </IntlProvider>
    );
}