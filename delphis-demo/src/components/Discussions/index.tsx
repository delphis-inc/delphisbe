import React, { useState } from 'react';
import classnames from 'classnames';
import './style.scss';
import DiscussionList from '../DiscussionList';

export interface Props {}

function Discussions(props: Props) {
    const [listSelected, setListSelected] = useState('public');

    return (
        <div className="discussions justify-center">
            <>    
                <div className="discussion-selector flex justify-left mt-6 ml-4 pb-6">
                    <div className={classnames({selector: true, selected: listSelected === 'public', 'mr-8': true})}>
                        <div onClick={() => setListSelected('public')} className="top-public">
                            Top Public Conversations
                        </div>
                    </div>
                    <div onClick={() => setListSelected('private')} className={classnames({selector: true, selected: listSelected === 'private'})}>
                        <div className="top-private">
                            Top Private Conversations
                        </div>
                    </div>
                </div>
                {DiscussionList({})}
            </>
        </div>
    );
}

export default Discussions;