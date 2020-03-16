import React, { ReactNode } from 'react';
import Header from '../Header';

export interface Props {
    children: ReactNode
}

function DelphisLayout(props: Props) {
    return (
        <div>
            <Header/>
            <div className="delphis-layout w-full flex justify-center">
                <div className="pl-2 pr-2 px-4 py-2 mx-6 w-full md:w-8/12 lg:w-8/12 xl:w-8/12 mx-6">
                    {props.children}
                </div>
            </div>
        </div>
    );
}

export default DelphisLayout;