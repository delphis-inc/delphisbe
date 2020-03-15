import React, { ReactNode } from 'react';
import Header from '../Header';

export interface Props {
    children: ReactNode
}

function DelphisLayout(props: Props) {
    return (
        <div className="delphis-layout">
            <Header/>
            {props.children}
        </div>
    );
}

export default DelphisLayout;